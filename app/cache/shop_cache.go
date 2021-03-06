/**
 * Copyright 2015 @ z3q.net.
 * name : shop_cache
 * author : jarryliu
 * date : 2016-06-03 17:22
 * description :
 * history :
 */
package cache

import (
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/format"
	"go2o/core/service/dps"
	"strconv"
)

// 设置商户站点配置
func GetShopDataKey(shopId int) string {
	return fmt.Sprintf("cache:online-shop:siteconf:%d", shopId)
}

func DelShopCache(merchantId int) {
	kvs := GetKVS()
	kvs.Del(GetValueMerchantCacheCK(merchantId))
	kvs.Del(GetMerchantSiteConfCK(merchantId))
}

// 根据主机头识别商店编号
func GetShopIdByHost(host string) (shopId int) {
	key := "go2o:cache:shop-host:" + host
	sto := GetKVS()
	var err error
	if shopId, err = sto.GetInt(key); err != nil || shopId <= 0 {
		_, shopId = dps.ShopService.GetShopIdByHost(host)
		if shopId > 0 {
			sto.SetExpire(key, shopId, DefaultMaxSeconds)
		}
	}
	return shopId
}

// 根据商城编号获取商户编号
func GetMchIdByShopId(shopId int) (mchId int) {
	key := "go2o:cache:mch-by-shop:" + strconv.Itoa(shopId)
	sto := GetKVS()
	var err error
	if mchId, err = sto.GetInt(key); err != nil || mchId <= 0 {
		mchId = dps.ShopService.GetMerchantId(shopId)
		if mchId > 0 {
			sto.SetExpire(key, mchId, DefaultMaxSeconds)
		} else {
			log.Println("[ Shop][ Exception] - shop", shopId, " no merchant bind!")
		}
	}
	return mchId
}

// 获取商城的数据
func GetOnlineShopData(shopId int) *shop.ShopDto {
	var v shop.ShopDto
	var sto gof.Storage = GetKVS()
	var key string = GetShopDataKey(shopId)
	if sto.Get(key, &v) != nil {
		mchId := GetMchIdByShopId(shopId)
		if v2 := dps.ShopService.GetShopData(mchId, shopId); v2 != nil {
			v3 := v2.Data.(shop.OnlineShop)
			v3.Logo = format.GetResUrl(v3.Logo)
			v2.Data = &v3
			if v2 != nil {
				sto.SetExpire(key, *v2, DefaultMaxSeconds)
			}
			return v2
		}
	}
	return &v
}
