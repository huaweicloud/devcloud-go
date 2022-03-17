/*
 * Copyright (c) 2013 The github.com/go-redis/redis Authors.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *    * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *    * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * 2021.11.15-Changed modify the constructor Client, split a file into
 * multiple small files, and extract duplicate code.
 * 			Huawei Technologies Co., Ltd.
 */

package commands_test

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	devspore "github.com/huaweicloud/devcloud-go/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sorted Sets commands", func() {
	ctx := context.TODO()
	var sortedSetsClient *devspore.DevsporeClient

	conf := configuration()
	BeforeEach(func() {
		sortedSetsClient = devspore.NewDevsporeClient(conf)
		Expect(sortedSetsClient.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(sortedSetsClient.Close()).NotTo(HaveOccurred())
	})

	Describe("sorted sets", func() {
		Describe("bzpopMax bepopMin", func() {
			BeforeEach(func() {
				err := sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{
					Score:  1,
					Member: "one",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{
					Score:  2,
					Member: "two",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{
					Score:  3,
					Member: "three",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
			})
			It("should BZPopMax", func() {
				member, err := sortedSetsClient.BZPopMax(ctx, 0, "zset1", "zset2").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(member).To(Equal(&redis.ZWithKey{
					Z: redis.Z{
						Score:  3,
						Member: "three",
					},
					Key: "zset1",
				}))
			})

			It("should BZPopMin", func() {
				member, err := sortedSetsClient.BZPopMin(ctx, 0, "zset1", "zset2").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(member).To(Equal(&redis.ZWithKey{
					Z: redis.Z{
						Score:  1,
						Member: "one",
					},
					Key: "zset1",
				}))
			})
		})

		var zsetPop = &redis.ZWithKey{
			Z: redis.Z{
				Member: "a",
				Score:  1,
			},
			Key: "zset",
		}

		It("should BZPopMax blocks", func() {
			started := make(chan bool)
			done := make(chan bool)
			go func() {
				defer GinkgoRecover()

				started <- true
				bZPopMax := sortedSetsClient.BZPopMax(ctx, 0, "zset")
				Expect(bZPopMax.Err()).NotTo(HaveOccurred())
				Expect(bZPopMax.Val()).To(Equal(zsetPop))
				done <- true
			}()
			<-started

			select {
			case _, _ = <-done:
				Fail("BZPopMax is not blocked")
			case <-time.After(time.Second):
				// ok
			}

			zAdd := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "a",
			})
			Expect(zAdd.Err()).NotTo(HaveOccurred())

			select {
			case _, _ = <-done:
				// ok
			case <-time.After(time.Second):
				Fail("BZPopMax is still blocked")
			}

			close(done)
			close(started)
		})

		It("should BZPopMax timeout", func() {
			val, err := sortedSetsClient.BZPopMax(ctx, time.Second, "zset1").Result()
			Expect(err).To(Equal(redis.Nil))
			Expect(val).To(BeNil())

			Expect(sortedSetsClient.Ping(ctx).Err()).NotTo(HaveOccurred())

			stats := sortedSetsClient.PoolStats()
			Expect(stats.Hits).To(Equal(uint32(2)))
			Expect(stats.Misses).To(Equal(uint32(1)))
			Expect(stats.Timeouts).To(Equal(uint32(0)))
		})

		It("should BZPopMin blocks", func() {
			started := make(chan bool)
			done := make(chan bool)
			go func() {
				defer GinkgoRecover()

				started <- true
				bZPopMin := sortedSetsClient.BZPopMin(ctx, 0, "zset")
				Expect(bZPopMin.Err()).NotTo(HaveOccurred())
				Expect(bZPopMin.Val()).To(Equal(zsetPop))
				done <- true
			}()
			<-started

			select {
			case _, _ = <-done:
				Fail("BZPopMin is not blocked")
			case <-time.After(time.Second):
				// ok
			}

			zAdd := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Member: "a",
				Score:  1,
			})
			Expect(zAdd.Err()).NotTo(HaveOccurred())

			select {
			case _, _ = <-done:
				// ok
			case <-time.After(time.Second):
				Fail("BZPopMin is still blocked")
			}

			close(started)
			close(done)
		})

		It("should BZPopMin timeout", func() {
			val, err := sortedSetsClient.BZPopMin(ctx, time.Second, "zset1").Result()
			Expect(err).To(Equal(redis.Nil))
			Expect(val).To(BeNil())

			Expect(sortedSetsClient.Ping(ctx).Err()).NotTo(HaveOccurred())

			stats := sortedSetsClient.PoolStats()
			Expect(stats.Hits).To(Equal(uint32(2)))
			Expect(stats.Misses).To(Equal(uint32(1)))
			Expect(stats.Timeouts).To(Equal(uint32(0)))
		})

		Describe("zadd, zadd bytes", func() {
			AfterEach(func() {
				vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(vals).To(Equal([]redis.Z{{
					Score:  1,
					Member: "one",
				}, {
					Score:  1,
					Member: "uno",
				}, {
					Score:  3,
					Member: "two",
				}}))
			})

			It("should ZAdd", func() {
				added, err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: "one",
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(1)))

				added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: "uno",
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(1)))

				added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  2,
					Member: "two",
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(1)))

				added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  3,
					Member: "two",
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(0)))
			})

			It("should ZAdd bytes", func() {
				added, err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: []byte("one"),
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(1)))

				added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: []byte("uno"),
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(1)))

				added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  2,
					Member: []byte("two"),
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(1)))

				added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  3,
					Member: []byte("two"),
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(0)))
			})
		})

		It("[redis6] should ZAddArgs", func() {
			// Test only the GT+LT options.
			added, err := sortedSetsClient.ZAddArgs(ctx, "zset", redis.ZAddArgs{
				GT:      true,
				Members: []redis.Z{{Score: 1, Member: "one"}},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))

			added, err = sortedSetsClient.ZAddArgs(ctx, "zset", redis.ZAddArgs{
				GT:      true,
				Members: []redis.Z{{Score: 2, Member: "one"}},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(0)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "one"}}))

			added, err = sortedSetsClient.ZAddArgs(ctx, "zset", redis.ZAddArgs{
				LT:      true,
				Members: []redis.Z{{Score: 1, Member: "one"}},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(0)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))
		})

		It("should ZAddNX", func() {
			added, err := sortedSetsClient.ZAddNX(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))

			added, err = sortedSetsClient.ZAddNX(ctx, "zset", &redis.Z{
				Score:  2,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(0)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))
		})

		It("should ZAddXX", func() {
			added, err := sortedSetsClient.ZAddXX(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(0)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(BeEmpty())

			added, err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(1)))

			added, err = sortedSetsClient.ZAddXX(ctx, "zset", &redis.Z{
				Score:  2,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(0)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "one"}}))
		})

		// TODO: remove in v9.
		It("should ZAddCh", func() {
			changed, err := sortedSetsClient.ZAddCh(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(changed).To(Equal(int64(1)))

			changed, err = sortedSetsClient.ZAddCh(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(changed).To(Equal(int64(0)))
		})

		// TODO: remove in v9.
		It("should ZAddNXCh", func() {
			changed, err := sortedSetsClient.ZAddNXCh(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(changed).To(Equal(int64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))

			changed, err = sortedSetsClient.ZAddNXCh(ctx, "zset", &redis.Z{
				Score:  2,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(changed).To(Equal(int64(0)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  1,
				Member: "one",
			}}))
		})

		// TODO: remove in v9.
		It("should ZAddXXCh", func() {
			changed, err := sortedSetsClient.ZAddXXCh(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(changed).To(Equal(int64(0)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(BeEmpty())

			added, err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(1)))

			changed, err = sortedSetsClient.ZAddXXCh(ctx, "zset", &redis.Z{
				Score:  2,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(changed).To(Equal(int64(1)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "one"}}))
		})

		// TODO: remove in v9.
		It("should ZIncr", func() {
			score, err := sortedSetsClient.ZIncr(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(score).To(Equal(float64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))

			score, err = sortedSetsClient.ZIncr(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(score).To(Equal(float64(2)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "one"}}))
		})

		// TODO: remove in v9.
		It("should ZIncrNX", func() {
			score, err := sortedSetsClient.ZIncrNX(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(score).To(Equal(float64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))

			score, err = sortedSetsClient.ZIncrNX(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).To(Equal(redis.Nil))
			Expect(score).To(Equal(float64(0)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))
		})

		// TODO: remove in v9.
		It("should ZIncrXX", func() {
			score, err := sortedSetsClient.ZIncrXX(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).To(Equal(redis.Nil))
			Expect(score).To(Equal(float64(0)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(BeEmpty())

			added, err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(added).To(Equal(int64(1)))

			score, err = sortedSetsClient.ZIncrXX(ctx, "zset", &redis.Z{
				Score:  1,
				Member: "one",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(score).To(Equal(float64(2)))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "one"}}))
		})

		It("should ZCard", func() {
			err := sortedSetsClient.ZAdd(ctx, "zcard", &redis.Z{
				Score:  1,
				Member: "one",
			}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zcard", &redis.Z{
				Score:  2,
				Member: "two",
			}).Err()
			Expect(err).NotTo(HaveOccurred())

			card, err := sortedSetsClient.ZCard(ctx, "zcard").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(card).To(Equal(int64(2)))
		})

		It("should ZIncrBy", func() {
			err := sortedSetsClient.ZAdd(ctx, "zincrby", &redis.Z{
				Score:  1,
				Member: "one",
			}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zincrby", &redis.Z{
				Score:  2,
				Member: "two",
			}).Err()
			Expect(err).NotTo(HaveOccurred())

			n, err := sortedSetsClient.ZIncrBy(ctx, "zincrby", 2, "one").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(float64(3)))

			val, err := sortedSetsClient.ZRangeWithScores(ctx, "zincrby", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal([]redis.Z{{
				Score:  2,
				Member: "two",
			}, {
				Score:  3,
				Member: "one",
			}}))
		})

		It("should ZInterStore", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{
				Score:  1,
				Member: "oneone",
			}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{
				Score:  2,
				Member: "twotwo",
			}).Err()
			Expect(err).NotTo(HaveOccurred())

			err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 1, Member: "oneone"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 2, Member: "twotwo"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset3", &redis.Z{Score: 3, Member: "twotwo"}).Err()
			Expect(err).NotTo(HaveOccurred())

			n, err := sortedSetsClient.ZInterStore(ctx, "out", &redis.ZStore{
				Keys:    []string{"zset1", "zset2"},
				Weights: []float64{2, 3},
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(2)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "out", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  5,
				Member: "oneone",
			}, {
				Score:  10,
				Member: "twotwo",
			}}))
		})

		It("[redis6] should ZMScore", func() {
			zmScore := sortedSetsClient.ZMScore(ctx, "zset", "one", "three")
			Expect(zmScore.Err()).NotTo(HaveOccurred())
			Expect(zmScore.Val()).To(HaveLen(2))
			Expect(zmScore.Val()[0]).To(Equal(float64(0)))

			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zmScore = sortedSetsClient.ZMScore(ctx, "zset", "one", "three")
			Expect(zmScore.Err()).NotTo(HaveOccurred())
			Expect(zmScore.Val()).To(HaveLen(2))
			Expect(zmScore.Val()[0]).To(Equal(float64(1)))

			zmScore = sortedSetsClient.ZMScore(ctx, "zset", "four")
			Expect(zmScore.Err()).NotTo(HaveOccurred())
			Expect(zmScore.Val()).To(HaveLen(1))

			zmScore = sortedSetsClient.ZMScore(ctx, "zset", "four", "one")
			Expect(zmScore.Err()).NotTo(HaveOccurred())
			Expect(zmScore.Val()).To(HaveLen(2))
		})

		Describe("zpopMin, zpopMax, zcount", func() {
			BeforeEach(func() {
				err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: "one",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  2,
					Member: "two",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  3,
					Member: "three",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
			})
			It("should ZCount", func() {
				count, err := sortedSetsClient.ZCount(ctx, "zset", "-inf", "+inf").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(3)))

				count, err = sortedSetsClient.ZCount(ctx, "zset", "(1", "3").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(2)))

				count, err = sortedSetsClient.ZLexCount(ctx, "zset", "-", "+").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(3)))
			})

			It("should ZPopMax", func() {
				members, err := sortedSetsClient.ZPopMax(ctx, "zset").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(members).To(Equal([]redis.Z{{
					Score:  3,
					Member: "three",
				}}))

				// adding back 3
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  3,
					Member: "three",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				members, err = sortedSetsClient.ZPopMax(ctx, "zset", 2).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(members).To(Equal([]redis.Z{{
					Score:  3,
					Member: "three",
				}, {
					Score:  2,
					Member: "two",
				}}))

				// adding back 2 & 3
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  3,
					Member: "three",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  2,
					Member: "two",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				members, err = sortedSetsClient.ZPopMax(ctx, "zset", 10).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(members).To(Equal([]redis.Z{{
					Score:  3,
					Member: "three",
				}, {
					Score:  2,
					Member: "two",
				}, {
					Score:  1,
					Member: "one",
				}}))
			})

			It("should ZPopMin", func() {
				members, err := sortedSetsClient.ZPopMin(ctx, "zset").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(members).To(Equal([]redis.Z{{
					Score:  1,
					Member: "one",
				}}))

				// adding back 1
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: "one",
				}).Err()
				Expect(err).NotTo(HaveOccurred())
				members, err = sortedSetsClient.ZPopMin(ctx, "zset", 2).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(members).To(Equal([]redis.Z{{
					Score:  1,
					Member: "one",
				}, {
					Score:  2,
					Member: "two",
				}}))

				// adding back 1 & 2
				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  1,
					Member: "one",
				}).Err()
				Expect(err).NotTo(HaveOccurred())

				err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
					Score:  2,
					Member: "two",
				}).Err()
				Expect(err).NotTo(HaveOccurred())

				members, err = sortedSetsClient.ZPopMin(ctx, "zset", 10).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(members).To(Equal([]redis.Z{{
					Score:  1,
					Member: "one",
				}, {
					Score:  2,
					Member: "two",
				}, {
					Score:  3,
					Member: "three",
				}}))
			})
		})

		It("should ZRange", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRange := sortedSetsClient.ZRange(ctx, "zset", 0, -1)
			Expect(zRange.Err()).NotTo(HaveOccurred())
			Expect(zRange.Val()).To(Equal([]string{"one", "two", "three"}))

			zRange = sortedSetsClient.ZRange(ctx, "zset", 2, 3)
			Expect(zRange.Err()).NotTo(HaveOccurred())
			Expect(zRange.Val()).To(Equal([]string{"three"}))

			zRange = sortedSetsClient.ZRange(ctx, "zset", -2, -1)
			Expect(zRange.Err()).NotTo(HaveOccurred())
			Expect(zRange.Val()).To(Equal([]string{"two", "three"}))
		})

		It("should ZRangeWithScores", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  1,
				Member: "one",
			}, {
				Score:  2,
				Member: "two",
			}, {
				Score:  3,
				Member: "three",
			}}))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", 2, 3).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 3, Member: "three"}}))

			vals, err = sortedSetsClient.ZRangeWithScores(ctx, "zset", -2, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  2,
				Member: "two",
			}, {
				Score:  3,
				Member: "three",
			}}))
		})

		It("should ZRangeByScore", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRangeByScore := sortedSetsClient.ZRangeByScore(ctx, "zset", &redis.ZRangeBy{
				Min: "-inf",
				Max: "+inf",
			})
			Expect(zRangeByScore.Err()).NotTo(HaveOccurred())
			Expect(zRangeByScore.Val()).To(Equal([]string{"one", "two", "three"}))

			zRangeByScore = sortedSetsClient.ZRangeByScore(ctx, "zset", &redis.ZRangeBy{
				Min: "1",
				Max: "2",
			})
			Expect(zRangeByScore.Err()).NotTo(HaveOccurred())
			Expect(zRangeByScore.Val()).To(Equal([]string{"one", "two"}))

			zRangeByScore = sortedSetsClient.ZRangeByScore(ctx, "zset", &redis.ZRangeBy{
				Min: "(1",
				Max: "2",
			})
			Expect(zRangeByScore.Err()).NotTo(HaveOccurred())
			Expect(zRangeByScore.Val()).To(Equal([]string{"two"}))

			zRangeByScore = sortedSetsClient.ZRangeByScore(ctx, "zset", &redis.ZRangeBy{
				Min: "(1",
				Max: "(2",
			})
			Expect(zRangeByScore.Err()).NotTo(HaveOccurred())
			Expect(zRangeByScore.Val()).To(Equal([]string{}))
		})

		It("should ZRangeByLex", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  0,
				Member: "a",
			}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  0,
				Member: "b",
			}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{
				Score:  0,
				Member: "c",
			}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRangeByLex := sortedSetsClient.ZRangeByLex(ctx, "zset", &redis.ZRangeBy{
				Min: "-",
				Max: "+",
			})
			Expect(zRangeByLex.Err()).NotTo(HaveOccurred())
			Expect(zRangeByLex.Val()).To(Equal([]string{"a", "b", "c"}))

			zRangeByLex = sortedSetsClient.ZRangeByLex(ctx, "zset", &redis.ZRangeBy{
				Min: "[a",
				Max: "[b",
			})
			Expect(zRangeByLex.Err()).NotTo(HaveOccurred())
			Expect(zRangeByLex.Val()).To(Equal([]string{"a", "b"}))

			zRangeByLex = sortedSetsClient.ZRangeByLex(ctx, "zset", &redis.ZRangeBy{
				Min: "(a",
				Max: "[b",
			})
			Expect(zRangeByLex.Err()).NotTo(HaveOccurred())
			Expect(zRangeByLex.Val()).To(Equal([]string{"b"}))

			zRangeByLex = sortedSetsClient.ZRangeByLex(ctx, "zset", &redis.ZRangeBy{
				Min: "(a",
				Max: "(b",
			})
			Expect(zRangeByLex.Err()).NotTo(HaveOccurred())
			Expect(zRangeByLex.Val()).To(Equal([]string{}))
		})

		It("should ZRangeByScoreWithScoresMap", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "oneone"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "twotwo"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "threethree"}).Err()
			Expect(err).NotTo(HaveOccurred())

			vals, err := sortedSetsClient.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
				Min: "-inf",
				Max: "+inf",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  1,
				Member: "oneone",
			}, {
				Score:  2,
				Member: "twotwo",
			}, {
				Score:  3,
				Member: "threethree",
			}}))

			vals, err = sortedSetsClient.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
				Min: "1",
				Max: "2",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  1,
				Member: "oneone",
			}, {
				Score:  2,
				Member: "twotwo",
			}}))

			vals, err = sortedSetsClient.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
				Min: "(1",
				Max: "2",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "twotwo"}}))

			vals, err = sortedSetsClient.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
				Min: "(1",
				Max: "(2",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{}))
		})

		Describe("zrangestore, zrangeargs", func() {
			BeforeEach(func() {
				added, err := sortedSetsClient.ZAddArgs(ctx, "zset", redis.ZAddArgs{
					Members: []redis.Z{
						{Score: 1, Member: "one"},
						{Score: 2, Member: "two"},
						{Score: 3, Member: "three"},
						{Score: 4, Member: "four"},
					},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(added).To(Equal(int64(4)))
			})

			It("[redis6] should ZRangeStore", func() {
				rangeStore, err := sortedSetsClient.ZRangeStore(ctx, "new-zset", redis.ZRangeArgs{
					Key:     "zset",
					Start:   1,
					Stop:    4,
					ByScore: true,
					Rev:     true,
					Offset:  1,
					Count:   2,
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(rangeStore).To(Equal(int64(2)))

				zRange, err := sortedSetsClient.ZRange(ctx, "new-zset", 0, -1).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(zRange).To(Equal([]string{"two", "three"}))
			})

			It("[redis6] should ZRangeArgs", func() {
				zRange, err := sortedSetsClient.ZRangeArgs(ctx, redis.ZRangeArgs{
					Key:     "zset",
					Start:   1,
					Stop:    4,
					ByScore: true,
					Rev:     true,
					Offset:  1,
					Count:   2,
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(zRange).To(Equal([]string{"three", "two"}))

				zRange, err = sortedSetsClient.ZRangeArgs(ctx, redis.ZRangeArgs{
					Key:    "zset",
					Start:  "-",
					Stop:   "+",
					ByLex:  true,
					Rev:    true,
					Offset: 2,
					Count:  2,
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(zRange).To(Equal([]string{"two", "one"}))

				zRange, err = sortedSetsClient.ZRangeArgs(ctx, redis.ZRangeArgs{
					Key:     "zset",
					Start:   "(1",
					Stop:    "(4",
					ByScore: true,
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(zRange).To(Equal([]string{"two", "three"}))

				// withScores.
				zSlice, err := sortedSetsClient.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
					Key:     "zset",
					Start:   1,
					Stop:    4,
					ByScore: true,
					Rev:     true,
					Offset:  1,
					Count:   2,
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(zSlice).To(Equal([]redis.Z{
					{Score: 3, Member: "three"},
					{Score: 2, Member: "two"},
				}))
			})
		})

		It("should ZRank", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRank := sortedSetsClient.ZRank(ctx, "zset", "three")
			Expect(zRank.Err()).NotTo(HaveOccurred())
			Expect(zRank.Val()).To(Equal(int64(2)))

			zRank = sortedSetsClient.ZRank(ctx, "zset", "four")
			Expect(zRank.Err()).To(Equal(redis.Nil))
			Expect(zRank.Val()).To(Equal(int64(0)))
		})

		It("should ZRem", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRem := sortedSetsClient.ZRem(ctx, "zset", "two")
			Expect(zRem.Err()).NotTo(HaveOccurred())
			Expect(zRem.Val()).To(Equal(int64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  1,
				Member: "one",
			}, {
				Score:  3,
				Member: "three",
			}}))
		})

		It("should ZRemRangeByRank", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRemRangeByRank := sortedSetsClient.ZRemRangeByRank(ctx, "zset", 0, 1)
			Expect(zRemRangeByRank.Err()).NotTo(HaveOccurred())
			Expect(zRemRangeByRank.Val()).To(Equal(int64(2)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  3,
				Member: "three",
			}}))
		})

		It("should ZRemRangeByScore", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRemRangeByScore := sortedSetsClient.ZRemRangeByScore(ctx, "zset", "-inf", "(2")
			Expect(zRemRangeByScore.Err()).NotTo(HaveOccurred())
			Expect(zRemRangeByScore.Val()).To(Equal(int64(1)))

			vals, err := sortedSetsClient.ZRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{
				Score:  2,
				Member: "two",
			}, {
				Score:  3,
				Member: "three",
			}}))
		})

		It("should ZRemRangeByLex", func() {
			zz := []*redis.Z{
				{Score: 0, Member: "aaaa"},
				{Score: 0, Member: "b"},
				{Score: 0, Member: "c"},
				{Score: 0, Member: "d"},
				{Score: 0, Member: "e"},
				{Score: 0, Member: "foo"},
				{Score: 0, Member: "zap"},
				{Score: 0, Member: "zip"},
				{Score: 0, Member: "ALPHA"},
				{Score: 0, Member: "alpha"},
			}
			for _, z := range zz {
				err := sortedSetsClient.ZAdd(ctx, "zset", z).Err()
				Expect(err).NotTo(HaveOccurred())
			}

			n, err := sortedSetsClient.ZRemRangeByLex(ctx, "zset", "[alpha", "[omega").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(6)))

			vals, err := sortedSetsClient.ZRange(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{"ALPHA", "aaaa", "zap", "zip"}))
		})

		It("should ZRevRange", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRevRange := sortedSetsClient.ZRevRange(ctx, "zset", 0, -1)
			Expect(zRevRange.Err()).NotTo(HaveOccurred())
			Expect(zRevRange.Val()).To(Equal([]string{"three", "two", "one"}))

			zRevRange = sortedSetsClient.ZRevRange(ctx, "zset", 2, 3)
			Expect(zRevRange.Err()).NotTo(HaveOccurred())
			Expect(zRevRange.Val()).To(Equal([]string{"one"}))

			zRevRange = sortedSetsClient.ZRevRange(ctx, "zset", -2, -1)
			Expect(zRevRange.Err()).NotTo(HaveOccurred())
			Expect(zRevRange.Val()).To(Equal([]string{"two", "one"}))
		})

		var zrangeRes = []redis.Z{{
			Score:  3,
			Member: "three",
		}, {
			Score:  2,
			Member: "two",
		}, {
			Score:  1,
			Member: "one",
		}}

		It("should ZRevRangeWithScoresMap", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			val, err := sortedSetsClient.ZRevRangeWithScores(ctx, "zset", 0, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(zrangeRes))

			val, err = sortedSetsClient.ZRevRangeWithScores(ctx, "zset", 2, 3).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal([]redis.Z{{Score: 1, Member: "one"}}))

			val, err = sortedSetsClient.ZRevRangeWithScores(ctx, "zset", -2, -1).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal([]redis.Z{{
				Score:  2,
				Member: "two",
			}, {
				Score:  1,
				Member: "one",
			}}))
		})

		It("should ZRevRangeByScore", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			vals, err := sortedSetsClient.ZRevRangeByScore(
				ctx, "zset", &redis.ZRangeBy{Max: "+inf", Min: "-inf"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{"three", "two", "one"}))

			vals, err = sortedSetsClient.ZRevRangeByScore(
				ctx, "zset", &redis.ZRangeBy{Max: "2", Min: "(1"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{"two"}))

			vals, err = sortedSetsClient.ZRevRangeByScore(
				ctx, "zset", &redis.ZRangeBy{Max: "(2", Min: "(1"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{}))
		})

		It("should ZRevRangeByLex", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 0, Member: "a"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 0, Member: "b"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 0, Member: "c"}).Err()
			Expect(err).NotTo(HaveOccurred())

			vals, err := sortedSetsClient.ZRevRangeByLex(
				ctx, "zset", &redis.ZRangeBy{Max: "+", Min: "-"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{"c", "b", "a"}))

			vals, err = sortedSetsClient.ZRevRangeByLex(
				ctx, "zset", &redis.ZRangeBy{Max: "[b", Min: "(a"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{"b"}))

			vals, err = sortedSetsClient.ZRevRangeByLex(
				ctx, "zset", &redis.ZRangeBy{Max: "(b", Min: "(a"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]string{}))
		})

		It("should ZRevRangeByScoreWithScoresMap", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			vals, err := sortedSetsClient.ZRevRangeByScoreWithScores(
				ctx, "zset", &redis.ZRangeBy{Max: "+inf", Min: "-inf"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal(zrangeRes))

			vals, err = sortedSetsClient.ZRevRangeByScoreWithScores(
				ctx, "zset", &redis.ZRangeBy{Max: "2", Min: "(1"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{{Score: 2, Member: "two"}}))

			vals, err = sortedSetsClient.ZRevRangeByScoreWithScores(
				ctx, "zset", &redis.ZRangeBy{Max: "(2", Min: "(1"}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(vals).To(Equal([]redis.Z{}))
		})

		It("should ZRevRank", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())

			zRevRank := sortedSetsClient.ZRevRank(ctx, "zset", "one")
			Expect(zRevRank.Err()).NotTo(HaveOccurred())
			Expect(zRevRank.Val()).To(Equal(int64(2)))

			zRevRank = sortedSetsClient.ZRevRank(ctx, "zset", "four")
			Expect(zRevRank.Err()).To(Equal(redis.Nil))
			Expect(zRevRank.Val()).To(Equal(int64(0)))
		})

		It("should ZScore", func() {
			zAdd := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1.001, Member: "one"})
			Expect(zAdd.Err()).NotTo(HaveOccurred())

			zScore := sortedSetsClient.ZScore(ctx, "zset", "one")
			Expect(zScore.Err()).NotTo(HaveOccurred())
			Expect(zScore.Val()).To(Equal(float64(1.001)))
		})

		It("[redis6] should ZUnion", func() {
			err := sortedSetsClient.ZAddArgs(ctx, "zset1", redis.ZAddArgs{
				Members: []redis.Z{
					{Score: 1, Member: "one"},
					{Score: 2, Member: "two"},
				},
			}).Err()
			Expect(err).NotTo(HaveOccurred())

			err = sortedSetsClient.ZAddArgs(ctx, "zset2", redis.ZAddArgs{
				Members: []redis.Z{
					{Score: 1, Member: "one"},
					{Score: 2, Member: "two"},
					{Score: 3, Member: "three"},
				},
			}).Err()
			Expect(err).NotTo(HaveOccurred())

			union, err := sortedSetsClient.ZUnion(ctx, redis.ZStore{
				Keys:      []string{"zset1", "zset2"},
				Weights:   []float64{2, 3},
				Aggregate: "sum",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(union).To(Equal([]string{"one", "three", "two"}))

			unionScores, err := sortedSetsClient.ZUnionWithScores(ctx, redis.ZStore{
				Keys:      []string{"zset1", "zset2"},
				Weights:   []float64{2, 3},
				Aggregate: "sum",
			}).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(unionScores).To(Equal([]redis.Z{
				{Score: 5, Member: "one"},
				{Score: 9, Member: "three"},
				{Score: 10, Member: "two"},
			}))
		})

		It("[redis6] should ZRandMember", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())

			v := sortedSetsClient.ZRandMember(ctx, "zset", 1, false)
			Expect(v.Err()).NotTo(HaveOccurred())
			Expect(v.Val()).To(Or(Equal([]string{"one"}), Equal([]string{"two"})))

			v = sortedSetsClient.ZRandMember(ctx, "zset", 0, false)
			Expect(v.Err()).NotTo(HaveOccurred())
			Expect(v.Val()).To(HaveLen(0))

			var slice []string
			err = sortedSetsClient.ZRandMember(ctx, "zset", 1, true).ScanSlice(&slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(slice).To(Or(Equal([]string{"one", "1"}), Equal([]string{"two", "2"})))
		})

		It("[redis6] should ZDiff", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())

			v, err := sortedSetsClient.ZDiff(ctx, "zset1", "zset2").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal([]string{"two", "three"}))
		})

		It("[redis6] should ZDiffWithScores", func() {
			err := sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 2, Member: "two"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 3, Member: "three"}).Err()
			Expect(err).NotTo(HaveOccurred())
			err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 1, Member: "one"}).Err()
			Expect(err).NotTo(HaveOccurred())

			v, err := sortedSetsClient.ZDiffWithScores(ctx, "zset1", "zset2").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal([]redis.Z{
				{
					Member: "two",
					Score:  2,
				},
				{
					Member: "three",
					Score:  3,
				},
			}))
		})

		Describe("zunionstore, zinter, zinterwithscore, zdiffstore", func() {
			BeforeEach(func() {
				err := sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 1, Member: "one"}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset1", &redis.Z{Score: 2, Member: "two"}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 1, Member: "one"}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 2, Member: "two"}).Err()
				Expect(err).NotTo(HaveOccurred())
				err = sortedSetsClient.ZAdd(ctx, "zset2", &redis.Z{Score: 3, Member: "three"}).Err()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should ZUnionStore", func() {
				n, err := sortedSetsClient.ZUnionStore(ctx, "out", &redis.ZStore{
					Keys:    []string{"zset1", "zset2"},
					Weights: []float64{2, 3},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(3)))

				val, err := sortedSetsClient.ZRangeWithScores(ctx, "out", 0, -1).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal([]redis.Z{{
					Score:  5,
					Member: "one",
				}, {
					Score:  9,
					Member: "three",
				}, {
					Score:  10,
					Member: "two",
				}}))
			})

			It("[redis6] should ZInter", func() {
				v, err := sortedSetsClient.ZInter(ctx, &redis.ZStore{
					Keys: []string{"zset1", "zset2"},
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(v).To(Equal([]string{"one", "two"}))
			})

			It("[redis6] should ZInterWithScores", func() {
				v, err := sortedSetsClient.ZInterWithScores(ctx, &redis.ZStore{
					Keys:      []string{"zset1", "zset2"},
					Weights:   []float64{2, 3},
					Aggregate: "Max",
				}).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(v).To(Equal([]redis.Z{
					{
						Member: "one",
						Score:  3,
					},
					{
						Member: "two",
						Score:  6,
					},
				}))
			})

			It("[redis6] should ZDiffStore", func() {
				v, err := sortedSetsClient.ZDiffStore(ctx, "out1", "zset1", "zset2").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(v).To(Equal(int64(0)))
				v, err = sortedSetsClient.ZDiffStore(ctx, "out1", "zset2", "zset1").Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(v).To(Equal(int64(1)))
				vals, err := sortedSetsClient.ZRangeWithScores(ctx, "out1", 0, -1).Result()
				Expect(err).NotTo(HaveOccurred())
				Expect(vals).To(Equal([]redis.Z{{
					Score:  3,
					Member: "three",
				}}))
			})
		})
	})

})
