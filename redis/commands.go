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
	"github.com/huaweicloud/devcloud-go/redis/strategy"
)

func (c *DevsporeClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Get(ctx, key)
}

func (c *DevsporeClient) Pipeline() redis.Pipeliner {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).Pipeline()
}

func (c *DevsporeClient) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).Pipelined(ctx, fn)
}

func (c *DevsporeClient) TxPipeline() redis.Pipeliner {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).TxPipeline()
}

func (c *DevsporeClient) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).TxPipelined(ctx, fn)
}

func (c *DevsporeClient) Command(ctx context.Context) *redis.CommandsInfoCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Command(ctx)
}

func (c *DevsporeClient) ClientGetName(ctx context.Context) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ClientGetName(ctx)
}

func (c *DevsporeClient) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Echo(ctx, message)
}

func (c *DevsporeClient) Ping(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Ping(ctx)
}

func (c *DevsporeClient) Quit(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Quit(ctx)
}

func (c *DevsporeClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Del(ctx, keys...)
}

func (c *DevsporeClient) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Unlink(ctx, keys...)
}

func (c *DevsporeClient) Dump(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Dump(ctx, key)
}

func (c *DevsporeClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Exists(ctx, keys...)
}

func (c *DevsporeClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Expire(ctx, key, expiration)
}

func (c *DevsporeClient) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ExpireAt(ctx, key, tm)
}

func (c *DevsporeClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Keys(ctx, pattern)
}

func (c *DevsporeClient) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Migrate(ctx, host, port, key, db, timeout)
}

func (c *DevsporeClient) Move(ctx context.Context, key string, db int) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Move(ctx, key, db)
}

func (c *DevsporeClient) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ObjectRefCount(ctx, key)
}

func (c *DevsporeClient) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ObjectEncoding(ctx, key)
}

func (c *DevsporeClient) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ObjectIdleTime(ctx, key)
}

func (c *DevsporeClient) Persist(ctx context.Context, key string) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Persist(ctx, key)
}

func (c *DevsporeClient) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).PExpire(ctx, key, expiration)
}

func (c *DevsporeClient) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).PExpireAt(ctx, key, tm)
}

func (c *DevsporeClient) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).PTTL(ctx, key)
}
func (c *DevsporeClient) RandomKey(ctx context.Context) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).RandomKey(ctx)
}

func (c *DevsporeClient) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Rename(ctx, key, newkey)
}

func (c *DevsporeClient) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).RenameNX(ctx, key, newkey)
}

func (c *DevsporeClient) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Restore(ctx, key, ttl, value)
}

func (c *DevsporeClient) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).RestoreReplace(ctx, key, ttl, value)
}

func (c *DevsporeClient) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Sort(ctx, key, sort)
}

func (c *DevsporeClient) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SortStore(ctx, key, store, sort)
}

func (c *DevsporeClient) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SortInterfaces(ctx, key, sort)
}

func (c *DevsporeClient) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Touch(ctx, keys...)
}

func (c *DevsporeClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).TTL(ctx, key)
}

func (c *DevsporeClient) Type(ctx context.Context, key string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Type(ctx, key)
}

func (c *DevsporeClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Append(ctx, key, value)
}

func (c *DevsporeClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Decr(ctx, key)
}

func (c *DevsporeClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).DecrBy(ctx, key, decrement)
}

func (c *DevsporeClient) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GetRange(ctx, key, start, end)
}

func (c *DevsporeClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GetSet(ctx, key, value)
}

func (c *DevsporeClient) GetEx(ctx context.Context, key string, expiration time.Duration) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GetEx(ctx, key, expiration)
}

func (c *DevsporeClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GetDel(ctx, key)
}

func (c *DevsporeClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Incr(ctx, key)
}

func (c *DevsporeClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).IncrBy(ctx, key, value)
}

func (c *DevsporeClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).IncrByFloat(ctx, key, value)
}

func (c *DevsporeClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).MGet(ctx, keys...)
}

func (c *DevsporeClient) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).MSet(ctx, values...)
}

func (c *DevsporeClient) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).MSetNX(ctx, values...)
}

func (c *DevsporeClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Set(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetArgs(ctx context.Context, key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SetArgs(ctx, key, value, a)
}

func (c *DevsporeClient) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SetEX(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SetNX(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SetXX(ctx, key, value, expiration)
}

func (c *DevsporeClient) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SetRange(ctx, key, offset, value)
}

func (c *DevsporeClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).StrLen(ctx, key)
}

func (c *DevsporeClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GetBit(ctx, key, offset)
}

func (c *DevsporeClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SetBit(ctx, key, offset, value)
}

func (c *DevsporeClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).BitCount(ctx, key, bitCount)
}

func (c *DevsporeClient) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BitOpAnd(ctx, destKey, keys...)
}

func (c *DevsporeClient) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BitOpOr(ctx, destKey, keys...)
}

func (c *DevsporeClient) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BitOpXor(ctx, destKey, keys...)
}

func (c *DevsporeClient) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BitOpNot(ctx, destKey, key)
}

func (c *DevsporeClient) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).BitPos(ctx, key, bit, pos...)
}

func (c *DevsporeClient) BitField(ctx context.Context, key string, args ...interface{}) *redis.IntSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BitField(ctx, key, args...)
}

func (c *DevsporeClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Scan(ctx, cursor, match, count)
}

func (c *DevsporeClient) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ScanType(ctx, cursor, match, count, keyType)
}

func (c *DevsporeClient) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SScan(ctx, key, cursor, match, count)
}

func (c *DevsporeClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HScan(ctx, key, cursor, match, count)
}

func (c *DevsporeClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZScan(ctx, key, cursor, match, count)
}

func (c *DevsporeClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).HDel(ctx, key, fields...)
}

func (c *DevsporeClient) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HExists(ctx, key, field)
}

func (c *DevsporeClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HGet(ctx, key, field)
}

func (c *DevsporeClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HGetAll(ctx, key)
}

func (c *DevsporeClient) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).HIncrBy(ctx, key, field, incr)
}

func (c *DevsporeClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).HIncrByFloat(ctx, key, field, incr)
}

func (c *DevsporeClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HKeys(ctx, key)
}

func (c *DevsporeClient) HLen(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HLen(ctx, key)
}

func (c *DevsporeClient) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HMGet(ctx, key, fields...)
}

func (c *DevsporeClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).HSet(ctx, key, values...)
}

func (c *DevsporeClient) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).HMSet(ctx, key, values...)
}

func (c *DevsporeClient) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).HSetNX(ctx, key, field, value)
}

func (c *DevsporeClient) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HVals(ctx, key)
}

func (c *DevsporeClient) HRandField(ctx context.Context, key string, count int, withValues bool) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).HRandField(ctx, key, count, withValues)
}

func (c *DevsporeClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BLPop(ctx, timeout, keys...)
}

func (c *DevsporeClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BRPop(ctx, timeout, keys...)
}

func (c *DevsporeClient) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BRPopLPush(ctx, source, destination, timeout)
}

func (c *DevsporeClient) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).LIndex(ctx, key, index)
}

func (c *DevsporeClient) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LInsert(ctx, key, op, pivot, value)
}

func (c *DevsporeClient) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LInsertBefore(ctx, key, pivot, value)
}

func (c *DevsporeClient) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LInsertAfter(ctx, key, pivot, value)
}

func (c *DevsporeClient) LLen(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).LLen(ctx, key)
}

func (c *DevsporeClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LPop(ctx, key)
}

func (c *DevsporeClient) LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LPopCount(ctx, key, count)
}

func (c *DevsporeClient) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).LPos(ctx, key, value, args)
}

func (c *DevsporeClient) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).LPosCount(ctx, key, value, count, args)
}

func (c *DevsporeClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LPush(ctx, key, values...)
}

func (c *DevsporeClient) LPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LPushX(ctx, key, values...)
}

func (c *DevsporeClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).LRange(ctx, key, start, stop)
}

func (c *DevsporeClient) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LRem(ctx, key, count, value)
}

func (c *DevsporeClient) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LSet(ctx, key, index, value)
}

func (c *DevsporeClient) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LTrim(ctx, key, start, stop)
}

func (c *DevsporeClient) RPop(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).RPop(ctx, key)
}

func (c *DevsporeClient) RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).RPopCount(ctx, key, count)
}

func (c *DevsporeClient) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).RPopLPush(ctx, source, destination)
}

func (c *DevsporeClient) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	return c.strategy.RouteClient(strategy.CommandTypeWrite).RPush(ctx, key, values...)
}

func (c *DevsporeClient) RPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).RPushX(ctx, key, values...)
}

func (c *DevsporeClient) LMove(ctx context.Context, source, destination, srcpos, destpos string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).LMove(ctx, source, destination, srcpos, destpos)
}

func (c *DevsporeClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SAdd(ctx, key, members...)
}

func (c *DevsporeClient) SCard(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SCard(ctx, key)
}

func (c *DevsporeClient) SDiff(ctx context.Context, key ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SDiff(ctx, key...)
}

func (c *DevsporeClient) SDiffStore(ctx context.Context, destination string, key ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SDiffStore(ctx, destination, key...)
}

func (c *DevsporeClient) SInter(ctx context.Context, key ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SInter(ctx, key...)
}

func (c *DevsporeClient) SInterStore(ctx context.Context, destination string, key ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SInterStore(ctx, destination, key...)
}

func (c *DevsporeClient) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SIsMember(ctx, key, member)
}

func (c *DevsporeClient) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SMIsMember(ctx, key, members...)
}

func (c *DevsporeClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SMembers(ctx, key)
}

func (c *DevsporeClient) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SMembersMap(ctx, key)
}

func (c *DevsporeClient) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SMove(ctx, source, destination, member)
}

func (c *DevsporeClient) SPop(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SPop(ctx, key)
}

func (c *DevsporeClient) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SPopN(ctx, key, count)
}

func (c *DevsporeClient) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SRandMember(ctx, key)
}

func (c *DevsporeClient) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SRandMemberN(ctx, key, count)
}

func (c *DevsporeClient) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SRem(ctx, key, members...)
}

func (c *DevsporeClient) SUnion(ctx context.Context, key ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SUnion(ctx, key...)
}

func (c *DevsporeClient) SUnionStore(ctx context.Context, destination string, key ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).SUnionStore(ctx, destination, key...)
}

func (c *DevsporeClient) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XAdd(ctx, a)
}

func (c *DevsporeClient) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XDel(ctx, stream, ids...)
}

func (c *DevsporeClient) XLen(ctx context.Context, stream string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XLen(ctx, stream)
}

func (c *DevsporeClient) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XRange(ctx, stream, start, stop)
}

func (c *DevsporeClient) XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XRangeN(ctx, stream, start, stop, count)
}

func (c *DevsporeClient) XRevRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XRevRange(ctx, stream, start, stop)
}

func (c *DevsporeClient) XRevRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XRevRangeN(ctx, stream, start, stop, count)
}

func (c *DevsporeClient) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XRead(ctx, a)
}

func (c *DevsporeClient) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XReadStreams(ctx, streams...)
}

func (c *DevsporeClient) XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XGroupCreate(ctx, stream, group, start)
}

func (c *DevsporeClient) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XGroupCreateMkStream(ctx, stream, group, start)
}

func (c *DevsporeClient) XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XGroupSetID(ctx, stream, group, start)
}

func (c *DevsporeClient) XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XGroupDestroy(ctx, stream, group)
}

func (c *DevsporeClient) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XGroupCreateConsumer(ctx, stream, group, consumer)
}

func (c *DevsporeClient) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XGroupDelConsumer(ctx, stream, group, consumer)
}

func (c *DevsporeClient) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XReadGroup(ctx, a)
}

func (c *DevsporeClient) XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XAck(ctx, stream, group, ids...)
}

func (c *DevsporeClient) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XPending(ctx, stream, group)
}

func (c *DevsporeClient) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XPendingExt(ctx, a)
}

func (c *DevsporeClient) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XClaim(ctx, a)
}

func (c *DevsporeClient) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XClaimJustID(ctx, a)
}

func (c *DevsporeClient) XAutoClaim(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XAutoClaim(ctx, a)
}

func (c *DevsporeClient) XAutoClaimJustID(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XAutoClaimJustID(ctx, a)
}

func (c *DevsporeClient) XTrim(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XTrim(ctx, key, maxLen)
}

func (c *DevsporeClient) XTrimApprox(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XTrimApprox(ctx, key, maxLen)
}

func (c *DevsporeClient) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XTrimMaxLen(ctx, key, maxLen)
}

func (c *DevsporeClient) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XTrimMaxLenApprox(ctx, key, maxLen, limit)
}

func (c *DevsporeClient) XTrimMinID(ctx context.Context, key string, minID string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XTrimMinID(ctx, key, minID)
}

func (c *DevsporeClient) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).XTrimMinIDApprox(ctx, key, minID, limit)
}

func (c *DevsporeClient) XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XInfoGroups(ctx, key)
}

func (c *DevsporeClient) XInfoStream(ctx context.Context, key string) *redis.XInfoStreamCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XInfoStream(ctx, key)
}

func (c *DevsporeClient) XInfoStreamFull(ctx context.Context, key string, count int) *redis.XInfoStreamFullCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XInfoStreamFull(ctx, key, count)
}

func (c *DevsporeClient) XInfoConsumers(ctx context.Context, key string, group string) *redis.XInfoConsumersCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).XInfoConsumers(ctx, key, group)
}

func (c *DevsporeClient) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BZPopMax(ctx, timeout, keys...)
}

func (c *DevsporeClient) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BZPopMin(ctx, timeout, keys...)
}

func (c *DevsporeClient) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAdd(ctx, key, members...)
}

func (c *DevsporeClient) ZAddNX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddNX(ctx, key, members...)
}

func (c *DevsporeClient) ZAddXX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddXX(ctx, key, members...)
}

func (c *DevsporeClient) ZAddCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddCh(ctx, key, members...)
}

func (c *DevsporeClient) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddNXCh(ctx, key, members...)
}

func (c *DevsporeClient) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddXXCh(ctx, key, members...)
}

func (c *DevsporeClient) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddArgs(ctx, key, args)
}

func (c *DevsporeClient) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZAddArgsIncr(ctx, key, args)
}

func (c *DevsporeClient) ZIncr(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZIncr(ctx, key, member)
}

func (c *DevsporeClient) ZIncrNX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZIncrNX(ctx, key, member)
}

func (c *DevsporeClient) ZIncrXX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZIncrXX(ctx, key, member)
}

func (c *DevsporeClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZCard(ctx, key)
}

func (c *DevsporeClient) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZCount(ctx, key, min, max)
}

func (c *DevsporeClient) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZLexCount(ctx, key, min, max)
}

func (c *DevsporeClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZIncrBy(ctx, key, increment, member)
}

func (c *DevsporeClient) ZInter(ctx context.Context, store *redis.ZStore) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZInter(ctx, store)
}

func (c *DevsporeClient) ZInterWithScores(ctx context.Context, store *redis.ZStore) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZInterWithScores(ctx, store)
}

func (c *DevsporeClient) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZInterStore(ctx, destination, store)
}

func (c *DevsporeClient) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZMScore(ctx, key, members...)
}

func (c *DevsporeClient) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZPopMax(ctx, key, count...)
}

func (c *DevsporeClient) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZPopMin(ctx, key, count...)
}

func (c *DevsporeClient) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRange(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRangeWithScores(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRangeByScore(ctx, key, opt)
}

func (c *DevsporeClient) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRangeByLex(ctx, key, opt)
}

func (c *DevsporeClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRangeByScoreWithScores(ctx, key, opt)
}

func (c *DevsporeClient) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRangeArgs(ctx, z)
}

func (c *DevsporeClient) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRangeArgsWithScores(ctx, z)
}

func (c *DevsporeClient) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZRangeStore(ctx, dst, z)
}

func (c *DevsporeClient) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRank(ctx, key, member)
}

func (c *DevsporeClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZRem(ctx, key, members...)
}

func (c *DevsporeClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZRemRangeByRank(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZRemRangeByScore(ctx, key, min, max)
}

func (c *DevsporeClient) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZRemRangeByLex(ctx, key, min, max)
}

func (c *DevsporeClient) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRevRange(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRevRangeWithScores(ctx, key, start, stop)
}

func (c *DevsporeClient) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRevRangeByScore(ctx, key, opt)
}

func (c *DevsporeClient) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRevRangeByLex(ctx, key, opt)
}

func (c *DevsporeClient) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (c *DevsporeClient) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRevRank(ctx, key, member)
}

func (c *DevsporeClient) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZScore(ctx, key, member)
}

func (c *DevsporeClient) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZUnionStore(ctx, dest, store)
}

func (c *DevsporeClient) ZUnion(ctx context.Context, store redis.ZStore) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZUnion(ctx, store)
}

func (c *DevsporeClient) ZUnionWithScores(ctx context.Context, store redis.ZStore) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZUnionWithScores(ctx, store)
}

func (c *DevsporeClient) ZRandMember(ctx context.Context, key string, count int, withScores bool) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZRandMember(ctx, key, count, withScores)
}

func (c *DevsporeClient) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZDiff(ctx, keys...)
}

func (c *DevsporeClient) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ZDiffWithScores(ctx, keys...)
}

func (c *DevsporeClient) ZDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ZDiffStore(ctx, destination, keys...)
}

func (c *DevsporeClient) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).PFAdd(ctx, key, els...)
}

func (c *DevsporeClient) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).PFCount(ctx, keys...)
}

func (c *DevsporeClient) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).PFMerge(ctx, dest, keys...)
}

func (c *DevsporeClient) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BgRewriteAOF(ctx)
}

func (c *DevsporeClient) BgSave(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).BgSave(ctx)
}

func (c *DevsporeClient) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ClientKill(ctx, ipPort)
}

func (c *DevsporeClient) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ClientKillByFilter(ctx, keys...)
}

func (c *DevsporeClient) ClientList(ctx context.Context) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ClientList(ctx)
}

func (c *DevsporeClient) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ClientPause(ctx, dur)
}

func (c *DevsporeClient) ClientID(ctx context.Context) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ClientID(ctx)
}

func (c *DevsporeClient) ConfigGet(ctx context.Context, parameter string) *redis.SliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ConfigGet(ctx, parameter)
}

func (c *DevsporeClient) ConfigResetStat(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ConfigResetStat(ctx)
}

func (c *DevsporeClient) ConfigSet(ctx context.Context, parameter, value string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ConfigSet(ctx, parameter, value)
}

func (c *DevsporeClient) ConfigRewrite(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ConfigRewrite(ctx)
}

func (c *DevsporeClient) DBSize(ctx context.Context) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).DBSize(ctx)
}

func (c *DevsporeClient) FlushAll(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).FlushAll(ctx)
}

func (c *DevsporeClient) FlushAllAsync(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).FlushAllAsync(ctx)
}

func (c *DevsporeClient) FlushDB(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).FlushDB(ctx)
}

func (c *DevsporeClient) FlushDBAsync(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).FlushDBAsync(ctx)
}

func (c *DevsporeClient) Info(ctx context.Context, section ...string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Info(ctx, section...)
}

func (c *DevsporeClient) LastSave(ctx context.Context) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).LastSave(ctx)
}

func (c *DevsporeClient) Save(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Save(ctx)
}

func (c *DevsporeClient) Shutdown(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Shutdown(ctx)
}

func (c *DevsporeClient) ShutdownSave(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ShutdownSave(ctx)
}

func (c *DevsporeClient) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).ShutdownNoSave(ctx)
}

func (c *DevsporeClient) SlaveOf(ctx context.Context, host, port string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).SlaveOf(ctx, host, port)
}

func (c *DevsporeClient) Time(ctx context.Context) *redis.TimeCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Time(ctx)
}

func (c *DevsporeClient) DebugObject(ctx context.Context, key string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).DebugObject(ctx, key)
}

func (c *DevsporeClient) ReadOnly(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ReadOnly(ctx)
}

func (c *DevsporeClient) ReadWrite(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).ReadWrite(ctx)
}

func (c *DevsporeClient) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).MemoryUsage(ctx, key, samples...)
}

func (c *DevsporeClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).Eval(ctx, script, keys, args...)
}

func (c *DevsporeClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).EvalSha(ctx, sha1, keys, args...)
}

func (c *DevsporeClient) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ScriptExists(ctx, hashes...)
}

func (c *DevsporeClient) ScriptFlush(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ScriptFlush(ctx)
}

func (c *DevsporeClient) ScriptKill(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ScriptKill(ctx)
}

func (c *DevsporeClient) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ScriptLoad(ctx, script)
}

func (c *DevsporeClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).Publish(ctx, channel, message)
}

func (c *DevsporeClient) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).PubSubChannels(ctx, pattern)
}

func (c *DevsporeClient) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).PubSubNumSub(ctx, channels...)
}

func (c *DevsporeClient) PubSubNumPat(ctx context.Context) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).PubSubNumPat(ctx)
}

func (c *DevsporeClient) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterSlots(ctx)
}

func (c *DevsporeClient) ClusterNodes(ctx context.Context) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterNodes(ctx)
}

func (c *DevsporeClient) ClusterMeet(ctx context.Context, host, port string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterMeet(ctx, host, port)
}

func (c *DevsporeClient) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterForget(ctx, nodeID)
}

func (c *DevsporeClient) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterReplicate(ctx, nodeID)
}

func (c *DevsporeClient) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterResetSoft(ctx)
}

func (c *DevsporeClient) ClusterResetHard(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterResetHard(ctx)
}

func (c *DevsporeClient) ClusterInfo(ctx context.Context) *redis.StringCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterInfo(ctx)
}

func (c *DevsporeClient) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterKeySlot(ctx, key)
}

func (c *DevsporeClient) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterGetKeysInSlot(ctx, slot, count)
}

func (c *DevsporeClient) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterCountFailureReports(ctx, nodeID)
}

func (c *DevsporeClient) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterCountKeysInSlot(ctx, slot)
}

func (c *DevsporeClient) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterDelSlots(ctx, slots...)
}

func (c *DevsporeClient) ClusterDelSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterDelSlotsRange(ctx, min, max)
}

func (c *DevsporeClient) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterSaveConfig(ctx)
}

func (c *DevsporeClient) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterSlaves(ctx, nodeID)
}

func (c *DevsporeClient) ClusterFailover(ctx context.Context) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterFailover(ctx)
}

func (c *DevsporeClient) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterAddSlots(ctx, slots...)
}

func (c *DevsporeClient) ClusterAddSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	return c.strategy.RouteClient(strategy.CommandTypeMulti).ClusterAddSlotsRange(ctx, min, max)
}

func (c *DevsporeClient) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GeoAdd(ctx, key, geoLocation...)
}

func (c *DevsporeClient) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoPos(ctx, key, members...)
}

func (c *DevsporeClient) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoRadius(ctx, key, longitude, latitude, query)
}

func (c *DevsporeClient) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GeoRadiusStore(ctx, key, longitude, latitude, query)
}

func (c *DevsporeClient) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoRadiusByMember(ctx, key, member, query)
}

func (c *DevsporeClient) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GeoRadiusByMemberStore(ctx, key, member, query)
}

func (c *DevsporeClient) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoSearch(ctx, key, q)
}

func (c *DevsporeClient) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) *redis.GeoSearchLocationCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoSearchLocation(ctx, key, q)
}

func (c *DevsporeClient) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) *redis.IntCmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).GeoSearchStore(ctx, key, store, q)
}

func (c *DevsporeClient) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoDist(ctx, key, member1, member2, unit)
}

func (c *DevsporeClient) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	return c.strategy.RouteClient(strategy.CommandTypeRead).GeoHash(ctx, key, members...)
}

func (c *DevsporeClient) PoolStats() *redis.PoolStats {
	return c.strategy.RouteClient(strategy.CommandTypeRead).PoolStats()
}

func (c *DevsporeClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.strategy.RouteClient(strategy.CommandTypeRead).Subscribe(ctx, channels...)
}

func (c *DevsporeClient) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.strategy.RouteClient(strategy.CommandTypeRead).PSubscribe(ctx, channels...)
}

func (c *DevsporeClient) Context() context.Context {
	return c.ctx
}

func (c *DevsporeClient) AddHook(hook redis.Hook) {
	c.strategy.RouteClient(strategy.CommandTypeWrite).AddHook(hook)
}

func (c *DevsporeClient) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return c.strategy.Watch(ctx, fn, keys...)
}

func (c *DevsporeClient) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Do(ctx, args...)
}

func (c *DevsporeClient) Process(ctx context.Context, cmd redis.Cmder) error {
	return c.strategy.RouteClient(strategy.CommandTypeWrite).Process(ctx, cmd)
}

func (c *DevsporeClient) SlowLogGet(ctx context.Context, num int64) *redis.SlowLogCmd {
	cmd := redis.NewSlowLogCmd(context.Background(), "slowlog", "get", num)
	_ = c.strategy.RouteClient(strategy.CommandTypeWrite).Process(ctx, cmd)
	return cmd
}

func (c *DevsporeClient) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "wait", numSlaves, int(timeout/time.Millisecond))
	_ = c.strategy.RouteClient(strategy.CommandTypeWrite).Process(ctx, cmd)
	return cmd
}

func (c *DevsporeClient) ClientUnblock(ctx context.Context, id int64) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "client", "unblock", id)
	_ = c.strategy.RouteClient(strategy.CommandTypeWrite).Process(ctx, cmd)
	return cmd
}

func (c *DevsporeClient) ClientUnblockWithError(ctx context.Context, id int64) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "client", "unblock", id, "error")
	_ = c.strategy.RouteClient(strategy.CommandTypeWrite).Process(ctx, cmd)
	return cmd
}
