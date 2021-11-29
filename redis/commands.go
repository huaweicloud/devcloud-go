/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2021.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type commandType int32

const (
	commandTypeRead commandType = iota
	commandTypeWrite
	commandTypeMulti
)

func (c *DevsporeClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).Get(ctx, key)
}

func (c *DevsporeClient) Pipeline() redis.Pipeliner {
	return c.getActualClient(commandTypeMulti).Pipeline()
}

func (c *DevsporeClient) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.getActualClient(commandTypeMulti).Pipelined(ctx, fn)
}

func (c *DevsporeClient) TxPipeline() redis.Pipeliner {
	return c.getActualClient(commandTypeMulti).TxPipeline()
}

func (c *DevsporeClient) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.getActualClient(commandTypeMulti).TxPipelined(ctx, fn)
}

func (c *DevsporeClient) Command(ctx context.Context) *redis.CommandsInfoCmd {
	return c.getActualClient(commandTypeRead).Command(ctx)
}

func (c *DevsporeClient) ClientGetName(ctx context.Context) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).ClientGetName(ctx)
}

func (c *DevsporeClient) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).Echo(ctx, message)
}

func (c *DevsporeClient) Ping(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeRead).Ping(ctx)
}

func (c *DevsporeClient) Quit(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeRead).Quit(ctx)
}

func (c *DevsporeClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).Del(ctx, keys...)
}

func (c *DevsporeClient) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).Unlink(ctx, keys...)
}

func (c *DevsporeClient) Dump(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).Dump(ctx, key)
}

func (c *DevsporeClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).Exists(ctx, keys...)
}

func (c *DevsporeClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).Expire(ctx, key, expiration)
}

func (c *DevsporeClient) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).ExpireAt(ctx, key, tm)
}

func (c *DevsporeClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).Keys(ctx, pattern)
}

func (c *DevsporeClient) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).Migrate(ctx, host, port, key, db, timeout)
}

func (c *DevsporeClient) Move(ctx context.Context, key string, db int) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).Move(ctx, key, db)
}

func (c *DevsporeClient) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ObjectRefCount(ctx, key)
}

func (c *DevsporeClient) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).ObjectEncoding(ctx, key)
}

func (c *DevsporeClient) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	return c.getActualClient(commandTypeRead).ObjectIdleTime(ctx, key)
}

func (c *DevsporeClient) Persist(ctx context.Context, key string) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).Persist(ctx, key)
}

func (c *DevsporeClient) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).PExpire(ctx, key, expiration)
}

func (c *DevsporeClient) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).PExpireAt(ctx, key, tm)
}

func (c *DevsporeClient) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	return c.getActualClient(commandTypeRead).PTTL(ctx, key)
}
func (c *DevsporeClient) RandomKey(ctx context.Context) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).RandomKey(ctx)
}

func (c *DevsporeClient) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).Rename(ctx, key, newkey)
}

func (c *DevsporeClient) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).RenameNX(ctx, key, newkey)
}

func (c *DevsporeClient) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).Restore(ctx, key, ttl, value)
}

func (c *DevsporeClient) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).RestoreReplace(ctx, key, ttl, value)
}

func (c *DevsporeClient) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).Sort(ctx, key, sort)
}

func (c *DevsporeClient) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).SortStore(ctx, key, store, sort)
}

func (c *DevsporeClient) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	return c.getActualClient(commandTypeRead).SortInterfaces(ctx, key, sort)
}

func (c *DevsporeClient) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).Touch(ctx, keys...)
}

func (c *DevsporeClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return c.getActualClient(commandTypeRead).TTL(ctx, key)
}

func (c *DevsporeClient) Type(ctx context.Context, key string) *redis.StatusCmd {
	return c.getActualClient(commandTypeRead).Type(ctx, key)
}

func (c *DevsporeClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).Append(ctx, key, value)
}

func (c *DevsporeClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).Decr(ctx, key)
}

func (c *DevsporeClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).DecrBy(ctx, key, decrement)
}

func (c *DevsporeClient) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).GetRange(ctx, key, start, end)
}

func (c *DevsporeClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).GetSet(ctx, key, value)
}

func (c *DevsporeClient) GetEx(ctx context.Context, key string, expiration time.Duration) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).GetEx(ctx, key, expiration)
}

func (c *DevsporeClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).GetDel(ctx, key)
}

func (c *DevsporeClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).Incr(ctx, key)
}

func (c *DevsporeClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).IncrBy(ctx, key, value)
}

func (c *DevsporeClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).IncrByFloat(ctx, key, value)
}

func (c *DevsporeClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return c.getActualClient(commandTypeRead).MGet(ctx, keys...)
}

func (c *DevsporeClient) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).MSet(ctx, values...)
}

func (c *DevsporeClient) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).MSetNX(ctx, values...)
}

func (c *DevsporeClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).Set(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetArgs(ctx context.Context, key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).SetArgs(ctx, key, value, a)
}

func (c *DevsporeClient) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).SetEX(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).SetNX(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).SetXX(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SetRange(ctx, key, offset, value)
}

func (c *DevsporeClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).StrLen(ctx, key)
}

func (c *DevsporeClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).GetBit(ctx, key, offset)
}

func (c *DevsporeClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SetBit(ctx, key, offset, value)
}

func (c *DevsporeClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).BitCount(ctx, key, bitCount)
}

func (c *DevsporeClient) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).BitOpAnd(ctx, destKey, keys...)
}

func (c *DevsporeClient) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).BitOpOr(ctx, destKey, keys...)
}

func (c *DevsporeClient) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).BitOpXor(ctx, destKey, keys...)
}

func (c *DevsporeClient) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).BitOpNot(ctx, destKey, key)
}

func (c *DevsporeClient) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).BitPos(ctx, key, bit, pos...)
}

func (c *DevsporeClient) BitField(ctx context.Context, key string, args ...interface{}) *redis.IntSliceCmd {
	return c.getActualClient(commandTypeWrite).BitField(ctx, key, args...)
}

func (c *DevsporeClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.getActualClient(commandTypeRead).Scan(ctx, cursor, match, count)
}

func (c *DevsporeClient) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return c.getActualClient(commandTypeRead).ScanType(ctx, cursor, match, count, keyType)
}

func (c *DevsporeClient) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.getActualClient(commandTypeRead).SScan(ctx, key, cursor, match, count)
}

func (c *DevsporeClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.getActualClient(commandTypeRead).HScan(ctx, key, cursor, match, count)
}

func (c *DevsporeClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.getActualClient(commandTypeRead).ZScan(ctx, key, cursor, match, count)
}

func (c *DevsporeClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).HDel(ctx, key, fields...)
}

func (c *DevsporeClient) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return c.getActualClient(commandTypeRead).HExists(ctx, key, field)
}

func (c *DevsporeClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).HGet(ctx, key, field)
}

func (c *DevsporeClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return c.getActualClient(commandTypeRead).HGetAll(ctx, key)
}

func (c *DevsporeClient) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).HIncrBy(ctx, key, field, incr)
}

func (c *DevsporeClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).HIncrByFloat(ctx, key, field, incr)
}

func (c *DevsporeClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).HKeys(ctx, key)
}

func (c *DevsporeClient) HLen(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).HLen(ctx, key)
}

func (c *DevsporeClient) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return c.getActualClient(commandTypeRead).HMGet(ctx, key, fields...)
}

func (c *DevsporeClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).HSet(ctx, key, values...)
}

func (c *DevsporeClient) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).HMSet(ctx, key, values...)
}

func (c *DevsporeClient) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).HSetNX(ctx, key, field, value)
}

func (c *DevsporeClient) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).HVals(ctx, key)
}

func (c *DevsporeClient) HRandField(ctx context.Context, key string, count int, withValues bool) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).HRandField(ctx, key, count, withValues)
}

func (c *DevsporeClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeWrite).BLPop(ctx, timeout, keys...)
}

func (c *DevsporeClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeWrite).BRPop(ctx, timeout, keys...)
}

func (c *DevsporeClient) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).BRPopLPush(ctx, source, destination, timeout)
}

func (c *DevsporeClient) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).LIndex(ctx, key, index)
}

func (c *DevsporeClient) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).LInsert(ctx, key, op, pivot, value)
}

func (c *DevsporeClient) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).LInsertBefore(ctx, key, pivot, value)
}

func (c *DevsporeClient) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).LInsertAfter(ctx, key, pivot, value)
}

func (c *DevsporeClient) LLen(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).LLen(ctx, key)
}

func (c *DevsporeClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).LPop(ctx, key)
}

func (c *DevsporeClient) LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeWrite).LPopCount(ctx, key, count)
}

func (c *DevsporeClient) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).LPos(ctx, key, value, args)
}

func (c *DevsporeClient) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {
	return c.getActualClient(commandTypeRead).LPosCount(ctx, key, value, count, args)
}

func (c *DevsporeClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).LPush(ctx, key, values...)
}

func (c *DevsporeClient) LPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).LPushX(ctx, key, values...)
}

func (c *DevsporeClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).LRange(ctx, key, start, stop)
}

func (c *DevsporeClient) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).LRem(ctx, key, count, value)
}

func (c *DevsporeClient) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).LSet(ctx, key, index, value)
}

func (c *DevsporeClient) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).LTrim(ctx, key, start, stop)
}

func (c *DevsporeClient) RPop(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).RPop(ctx, key)
}

func (c *DevsporeClient) RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeWrite).RPopCount(ctx, key, count)
}

func (c *DevsporeClient) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).RPopLPush(ctx, source, destination)
}

func (c *DevsporeClient) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).RPush(ctx, key, values...)
}

func (c *DevsporeClient) RPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).RPushX(ctx, key, values...)
}

func (c *DevsporeClient) LMove(ctx context.Context, source, destination, srcpos, destpos string) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).LMove(ctx, source, destination, srcpos, destpos)
}

func (c *DevsporeClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SAdd(ctx, key, members...)
}

func (c *DevsporeClient) SCard(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).SCard(ctx, key)
}

func (c *DevsporeClient) SDiff(ctx context.Context, key ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).SDiff(ctx, key...)
}

func (c *DevsporeClient) SDiffStore(ctx context.Context, destination string, key ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SDiffStore(ctx, destination, key...)
}

func (c *DevsporeClient) SInter(ctx context.Context, key ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).SInter(ctx, key...)
}

func (c *DevsporeClient) SInterStore(ctx context.Context, destination string, key ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SInterStore(ctx, destination, key...)
}

func (c *DevsporeClient) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return c.getActualClient(commandTypeRead).SIsMember(ctx, key, member)
}

func (c *DevsporeClient) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	return c.getActualClient(commandTypeRead).SMIsMember(ctx, key, members...)
}

func (c *DevsporeClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).SMembers(ctx, key)
}

func (c *DevsporeClient) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	return c.getActualClient(commandTypeRead).SMembersMap(ctx, key)
}

func (c *DevsporeClient) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return c.getActualClient(commandTypeWrite).SMove(ctx, source, destination, member)
}

func (c *DevsporeClient) SPop(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).SPop(ctx, key)
}

func (c *DevsporeClient) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeWrite).SPopN(ctx, key, count)
}

func (c *DevsporeClient) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).SRandMember(ctx, key)
}

func (c *DevsporeClient) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).SRandMemberN(ctx, key, count)
}

func (c *DevsporeClient) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SRem(ctx, key, members...)
}

func (c *DevsporeClient) SUnion(ctx context.Context, key ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).SUnion(ctx, key...)
}

func (c *DevsporeClient) SUnionStore(ctx context.Context, destination string, key ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).SUnionStore(ctx, destination, key...)
}

func (c *DevsporeClient) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	return c.getActualClient(commandTypeWrite).XAdd(ctx, a)
}

func (c *DevsporeClient) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XDel(ctx, stream, ids...)
}

func (c *DevsporeClient) XLen(ctx context.Context, stream string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).XLen(ctx, stream)
}

func (c *DevsporeClient) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	return c.getActualClient(commandTypeRead).XRange(ctx, stream, start, stop)
}

func (c *DevsporeClient) XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return c.getActualClient(commandTypeRead).XRangeN(ctx, stream, start, stop, count)
}

func (c *DevsporeClient) XRevRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	return c.getActualClient(commandTypeRead).XRevRange(ctx, stream, start, stop)
}

func (c *DevsporeClient) XRevRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return c.getActualClient(commandTypeRead).XRevRangeN(ctx, stream, start, stop, count)
}

func (c *DevsporeClient) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return c.getActualClient(commandTypeRead).XRead(ctx, a)
}

func (c *DevsporeClient) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	return c.getActualClient(commandTypeRead).XReadStreams(ctx, streams...)
}

func (c *DevsporeClient) XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).XGroupCreate(ctx, stream, group, start)
}

func (c *DevsporeClient) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).XGroupCreateMkStream(ctx, stream, group, start)
}

func (c *DevsporeClient) XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).XGroupSetID(ctx, stream, group, start)
}

func (c *DevsporeClient) XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XGroupDestroy(ctx, stream, group)
}

func (c *DevsporeClient) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XGroupCreateConsumer(ctx, stream, group, consumer)
}

func (c *DevsporeClient) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XGroupDelConsumer(ctx, stream, group, consumer)
}

func (c *DevsporeClient) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	return c.getActualClient(commandTypeRead).XReadGroup(ctx, a)
}

func (c *DevsporeClient) XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XAck(ctx, stream, group, ids...)
}

func (c *DevsporeClient) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {
	return c.getActualClient(commandTypeWrite).XPending(ctx, stream, group)
}

func (c *DevsporeClient) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	return c.getActualClient(commandTypeWrite).XPendingExt(ctx, a)
}

func (c *DevsporeClient) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	return c.getActualClient(commandTypeWrite).XClaim(ctx, a)
}

func (c *DevsporeClient) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeWrite).XClaimJustID(ctx, a)
}

func (c *DevsporeClient) XAutoClaim(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {
	return c.getActualClient(commandTypeWrite).XAutoClaim(ctx, a)
}

func (c *DevsporeClient) XAutoClaimJustID(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {
	return c.getActualClient(commandTypeWrite).XAutoClaimJustID(ctx, a)
}

func (c *DevsporeClient) XTrim(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XTrim(ctx, key, maxLen)
}

func (c *DevsporeClient) XTrimApprox(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XTrimApprox(ctx, key, maxLen)
}

func (c *DevsporeClient) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XTrimMaxLen(ctx, key, maxLen)
}

func (c *DevsporeClient) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XTrimMaxLenApprox(ctx, key, maxLen, limit)
}

func (c *DevsporeClient) XTrimMinID(ctx context.Context, key string, minID string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XTrimMinID(ctx, key, minID)
}

func (c *DevsporeClient) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).XTrimMinIDApprox(ctx, key, minID, limit)
}

func (c *DevsporeClient) XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd {
	return c.getActualClient(commandTypeRead).XInfoGroups(ctx, key)
}

func (c *DevsporeClient) XInfoStream(ctx context.Context, key string) *redis.XInfoStreamCmd {
	return c.getActualClient(commandTypeRead).XInfoStream(ctx, key)
}

func (c *DevsporeClient) XInfoStreamFull(ctx context.Context, key string, count int) *redis.XInfoStreamFullCmd {
	return c.getActualClient(commandTypeRead).XInfoStreamFull(ctx, key, count)
}

func (c *DevsporeClient) XInfoConsumers(ctx context.Context, key string, group string) *redis.XInfoConsumersCmd {
	return c.getActualClient(commandTypeRead).XInfoConsumers(ctx, key, group)
}

func (c *DevsporeClient) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return c.getActualClient(commandTypeWrite).BZPopMax(ctx, timeout, keys...)
}

func (c *DevsporeClient) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return c.getActualClient(commandTypeWrite).BZPopMin(ctx, timeout, keys...)
}

func (c *DevsporeClient) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAdd(ctx, key, members...)
}

func (c *DevsporeClient) ZAddNX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAddNX(ctx, key, members...)
}

func (c *DevsporeClient) ZAddXX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAddXX(ctx, key, members...)
}

func (c *DevsporeClient) ZAddCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAddCh(ctx, key, members...)
}

func (c *DevsporeClient) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAddNXCh(ctx, key, members...)
}

func (c *DevsporeClient) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAddXXCh(ctx, key, members...)
}

func (c *DevsporeClient) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZAddArgs(ctx, key, args)
}

func (c *DevsporeClient) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).ZAddArgsIncr(ctx, key, args)
}

func (c *DevsporeClient) ZIncr(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).ZIncr(ctx, key, member)
}

func (c *DevsporeClient) ZIncrNX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).ZIncrNX(ctx, key, member)
}

func (c *DevsporeClient) ZIncrXX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).ZIncrXX(ctx, key, member)
}

func (c *DevsporeClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ZCard(ctx, key)
}

func (c *DevsporeClient) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ZCount(ctx, key, min, max)
}

func (c *DevsporeClient) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ZLexCount(ctx, key, min, max)
}

func (c *DevsporeClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	return c.getActualClient(commandTypeWrite).ZIncrBy(ctx, key, increment, member)
}

func (c *DevsporeClient) ZInter(ctx context.Context, store *redis.ZStore) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZInter(ctx, store)
}

func (c *DevsporeClient) ZInterWithScores(ctx context.Context, store *redis.ZStore) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZInterWithScores(ctx, store)
}

func (c *DevsporeClient) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZInterStore(ctx, destination, store)
}

func (c *DevsporeClient) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {
	return c.getActualClient(commandTypeRead).ZMScore(ctx, key, members...)
}

func (c *DevsporeClient) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeWrite).ZPopMax(ctx, key, count...)
}

func (c *DevsporeClient) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeWrite).ZPopMin(ctx, key, count...)
}

func (c *DevsporeClient) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRange(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZRangeWithScores(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRangeByScore(ctx, key, opt)
}

func (c *DevsporeClient) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRangeByLex(ctx, key, opt)
}

func (c *DevsporeClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZRangeByScoreWithScores(ctx, key, opt)
}

func (c *DevsporeClient) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRangeArgs(ctx, z)
}

func (c *DevsporeClient) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZRangeArgsWithScores(ctx, z)
}

func (c *DevsporeClient) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZRangeStore(ctx, dst, z)
}

func (c *DevsporeClient) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ZRank(ctx, key, member)
}

func (c *DevsporeClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZRem(ctx, key, members...)
}

func (c *DevsporeClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZRemRangeByRank(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZRemRangeByScore(ctx, key, min, max)
}

func (c *DevsporeClient) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZRemRangeByLex(ctx, key, min, max)
}

func (c *DevsporeClient) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRevRange(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZRevRangeWithScores(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRevRangeByScore(ctx, key, opt)
}

func (c *DevsporeClient) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRevRangeByLex(ctx, key, opt)
}

func (c *DevsporeClient) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (c *DevsporeClient) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ZRevRank(ctx, key, member)
}

func (c *DevsporeClient) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	return c.getActualClient(commandTypeRead).ZScore(ctx, key, member)
}

func (c *DevsporeClient) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZUnionStore(ctx, dest, store)
}

func (c *DevsporeClient) ZUnion(ctx context.Context, store redis.ZStore) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZUnion(ctx, store)
}

func (c *DevsporeClient) ZUnionWithScores(ctx context.Context, store redis.ZStore) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZUnionWithScores(ctx, store)
}

func (c *DevsporeClient) ZRandMember(ctx context.Context, key string, count int, withScores bool) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZRandMember(ctx, key, count, withScores)
}

func (c *DevsporeClient) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).ZDiff(ctx, keys...)
}

func (c *DevsporeClient) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {
	return c.getActualClient(commandTypeRead).ZDiffWithScores(ctx, keys...)
}

func (c *DevsporeClient) ZDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ZDiffStore(ctx, destination, keys...)
}

func (c *DevsporeClient) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).PFAdd(ctx, key, els...)
}

func (c *DevsporeClient) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).PFCount(ctx, keys...)
}

func (c *DevsporeClient) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).PFMerge(ctx, dest, keys...)
}

func (c *DevsporeClient) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).BgRewriteAOF(ctx)
}

func (c *DevsporeClient) BgSave(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).BgSave(ctx)
}

func (c *DevsporeClient) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).ClientKill(ctx, ipPort)
}

func (c *DevsporeClient) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).ClientKillByFilter(ctx, keys...)
}

func (c *DevsporeClient) ClientList(ctx context.Context) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).ClientList(ctx)
}

func (c *DevsporeClient) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {
	return c.getActualClient(commandTypeRead).ClientPause(ctx, dur)
}

func (c *DevsporeClient) ClientID(ctx context.Context) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).ClientID(ctx)
}

func (c *DevsporeClient) ConfigGet(ctx context.Context, parameter string) *redis.SliceCmd {
	return c.getActualClient(commandTypeRead).ConfigGet(ctx, parameter)
}

func (c *DevsporeClient) ConfigResetStat(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).ConfigResetStat(ctx)
}

func (c *DevsporeClient) ConfigSet(ctx context.Context, parameter, value string) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).ConfigSet(ctx, parameter, value)
}

func (c *DevsporeClient) ConfigRewrite(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).ConfigRewrite(ctx)
}

func (c *DevsporeClient) DBSize(ctx context.Context) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).DBSize(ctx)
}

func (c *DevsporeClient) FlushAll(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).FlushAll(ctx)
}

func (c *DevsporeClient) FlushAllAsync(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).FlushAllAsync(ctx)
}

func (c *DevsporeClient) FlushDB(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).FlushDB(ctx)
}

func (c *DevsporeClient) FlushDBAsync(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).FlushDBAsync(ctx)
}

func (c *DevsporeClient) Info(ctx context.Context, section ...string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).Info(ctx, section...)
}

func (c *DevsporeClient) LastSave(ctx context.Context) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).LastSave(ctx)
}

func (c *DevsporeClient) Save(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).Save(ctx)
}

func (c *DevsporeClient) Shutdown(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).Shutdown(ctx)
}

func (c *DevsporeClient) ShutdownSave(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).ShutdownSave(ctx)
}

func (c *DevsporeClient) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeWrite).ShutdownNoSave(ctx)
}

func (c *DevsporeClient) SlaveOf(ctx context.Context, host, port string) *redis.StatusCmd {
	return c.getActualClient(commandTypeRead).SlaveOf(ctx, host, port)
}

func (c *DevsporeClient) Time(ctx context.Context) *redis.TimeCmd {
	return c.getActualClient(commandTypeRead).Time(ctx)
}

func (c *DevsporeClient) DebugObject(ctx context.Context, key string) *redis.StringCmd {
	return c.getActualClient(commandTypeRead).DebugObject(ctx, key)
}

func (c *DevsporeClient) ReadOnly(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeRead).ReadOnly(ctx)
}

func (c *DevsporeClient) ReadWrite(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeRead).ReadWrite(ctx)
}

func (c *DevsporeClient) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	return c.getActualClient(commandTypeRead).MemoryUsage(ctx, key, samples...)
}

func (c *DevsporeClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return c.getActualClient(commandTypeMulti).Eval(ctx, script, keys, args...)
}

func (c *DevsporeClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return c.getActualClient(commandTypeMulti).EvalSha(ctx, sha1, keys, args...)
}

func (c *DevsporeClient) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return c.getActualClient(commandTypeMulti).ScriptExists(ctx, hashes...)
}

func (c *DevsporeClient) ScriptFlush(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ScriptFlush(ctx)
}

func (c *DevsporeClient) ScriptKill(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ScriptKill(ctx)
}

func (c *DevsporeClient) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	return c.getActualClient(commandTypeMulti).ScriptLoad(ctx, script)
}

func (c *DevsporeClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return c.getActualClient(commandTypeMulti).Publish(ctx, channel, message)
}

func (c *DevsporeClient) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeMulti).PubSubChannels(ctx, pattern)
}

func (c *DevsporeClient) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	return c.getActualClient(commandTypeMulti).PubSubNumSub(ctx, channels...)
}

func (c *DevsporeClient) PubSubNumPat(ctx context.Context) *redis.IntCmd {
	return c.getActualClient(commandTypeMulti).PubSubNumPat(ctx)
}

func (c *DevsporeClient) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {
	return c.getActualClient(commandTypeMulti).ClusterSlots(ctx)
}

func (c *DevsporeClient) ClusterNodes(ctx context.Context) *redis.StringCmd {
	return c.getActualClient(commandTypeMulti).ClusterNodes(ctx)
}

func (c *DevsporeClient) ClusterMeet(ctx context.Context, host, port string) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterMeet(ctx, host, port)
}

func (c *DevsporeClient) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterForget(ctx, nodeID)
}

func (c *DevsporeClient) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterReplicate(ctx, nodeID)
}

func (c *DevsporeClient) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterResetSoft(ctx)
}

func (c *DevsporeClient) ClusterResetHard(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterResetHard(ctx)
}

func (c *DevsporeClient) ClusterInfo(ctx context.Context) *redis.StringCmd {
	return c.getActualClient(commandTypeMulti).ClusterInfo(ctx)
}

func (c *DevsporeClient) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {
	return c.getActualClient(commandTypeMulti).ClusterKeySlot(ctx, key)
}

func (c *DevsporeClient) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeMulti).ClusterGetKeysInSlot(ctx, slot, count)
}

func (c *DevsporeClient) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	return c.getActualClient(commandTypeMulti).ClusterCountFailureReports(ctx, nodeID)
}

func (c *DevsporeClient) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {
	return c.getActualClient(commandTypeMulti).ClusterCountKeysInSlot(ctx, slot)
}

func (c *DevsporeClient) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterDelSlots(ctx, slots...)
}

func (c *DevsporeClient) ClusterDelSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterDelSlotsRange(ctx, min, max)
}

func (c *DevsporeClient) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterSaveConfig(ctx)
}

func (c *DevsporeClient) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeMulti).ClusterSlaves(ctx, nodeID)
}

func (c *DevsporeClient) ClusterFailover(ctx context.Context) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterFailover(ctx)
}

func (c *DevsporeClient) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterAddSlots(ctx, slots...)
}

func (c *DevsporeClient) ClusterAddSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	return c.getActualClient(commandTypeMulti).ClusterAddSlotsRange(ctx, min, max)
}

func (c *DevsporeClient) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).GeoAdd(ctx, key, geoLocation...)
}

func (c *DevsporeClient) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	return c.getActualClient(commandTypeRead).GeoPos(ctx, key, members...)
}

func (c *DevsporeClient) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return c.getActualClient(commandTypeRead).GeoRadius(ctx, key, longitude, latitude, query)
}

func (c *DevsporeClient) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).GeoRadiusStore(ctx, key, longitude, latitude, query)
}

func (c *DevsporeClient) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return c.getActualClient(commandTypeRead).GeoRadiusByMember(ctx, key, member, query)
}

func (c *DevsporeClient) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).GeoRadiusByMemberStore(ctx, key, member, query)
}

func (c *DevsporeClient) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).GeoSearch(ctx, key, q)
}

func (c *DevsporeClient) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) *redis.GeoSearchLocationCmd {
	return c.getActualClient(commandTypeRead).GeoSearchLocation(ctx, key, q)
}

func (c *DevsporeClient) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) *redis.IntCmd {
	return c.getActualClient(commandTypeWrite).GeoSearchStore(ctx, key, store, q)
}

func (c *DevsporeClient) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {
	return c.getActualClient(commandTypeRead).GeoDist(ctx, key, member1, member2, unit)
}

func (c *DevsporeClient) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	return c.getActualClient(commandTypeRead).GeoHash(ctx, key, members...)
}

func (c *DevsporeClient) PoolStats() *redis.PoolStats {
	return c.getActualClient(commandTypeRead).PoolStats()
}

func (c *DevsporeClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.getActualClient(commandTypeRead).Subscribe(ctx, channels...)
}

func (c *DevsporeClient) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.getActualClient(commandTypeRead).PSubscribe(ctx, channels...)
}

func (c *DevsporeClient) Context() context.Context {
	return c.ctx
}

func (c *DevsporeClient) AddHook(hook redis.Hook) {
	c.getActualClient(commandTypeWrite).AddHook(hook)
}

func (c *DevsporeClient) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return c.getActualClient(commandTypeWrite).Watch(ctx, fn, keys...)
}

func (c *DevsporeClient) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return c.getActualClient(commandTypeWrite).Do(ctx, args...)
}

func (c *DevsporeClient) Process(ctx context.Context, cmd redis.Cmder) error {
	return c.getActualClient(commandTypeWrite).Process(ctx, cmd)
}

func (c *DevsporeClient) SlowLogGet(ctx context.Context, num int64) *redis.SlowLogCmd {
	cmd := redis.NewSlowLogCmd(context.Background(), "slowlog", "get", num)
	_ = c.getActualClient(commandTypeWrite).Process(ctx, cmd)
	return cmd
}

func (c *DevsporeClient) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "wait", numSlaves, int(timeout/time.Millisecond))
	_ = c.getActualClient(commandTypeWrite).Process(ctx, cmd)
	return cmd
}

func (c *DevsporeClient) ClientUnblock(ctx context.Context, id int64) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "client", "unblock", id)
	_ = c.getActualClient(commandTypeWrite).Process(ctx, cmd)
	return cmd
}

func (c *DevsporeClient) ClientUnblockWithError(ctx context.Context, id int64) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "client", "unblock", id, "error")
	_ = c.getActualClient(commandTypeWrite).Process(ctx, cmd)
	return cmd
}
