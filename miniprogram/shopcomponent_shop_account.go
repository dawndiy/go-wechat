package miniprogram

import (
	"context"
)

// ShopCategoryDetail 商品类目详情
type ShopCategoryDetail struct {
	ShopCategory             `json:",inline"`
	Qualification            string `json:"qualification"`
	QualificationType        int    `json:"qualification_type"`
	ProductQualification     string `json:"product_qualification"`
	ProductQualificationType int    `json:"product_qualification_type"`
}

// Category 获取商品类目
//
// 获取所有三级类目及其资质相关信息 注意：该接口拉到的是【全量】三级类目数据，数据回包大小约为2MB。
// 所以请商家自己做好缓存，不要频繁调用（有严格的频率限制），该类目数据不会频率变动，推荐商家每天调用一次更新商家自身缓存
//
// 若该类目资质必填，则新增商品前，必须先通过该类目资质申请接口进行资质申请;
// 若该类目资质不需要，则该类目自动拥有，无需申请，如依然调用，会报错1050011； 若该商品资质必填，则新增商品时，带上商品资质字段。
// 接入类目审核回调，才可获取审核结果。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/cat/get_children_cateogry.html
func (s *ShopComponentShopService) Category(ctx context.Context) ([]ShopCategoryDetail, error) {
	u, err := s.client.apiURL(ctx, "shop/cat/get", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return nil, err
	}
	var data struct {
		ThirdCatList []ShopCategoryDetail `json:"third_cat_list"`
	}
	_, err = s.client.Do(req, &data)
	return data.ThirdCatList, err
}

// ShopCategoryDetail 商品类目详情
type ShopCategory struct {
	FirstCatID    int    `json:"first_cat_id"`
	FirstCatName  string `json:"first_cat_name"`
	SecondCatID   int    `json:"second_cat_id"`
	SecondCatName string `json:"second_cat_name"`
	ThirdCatID    int    `json:"third_cat_id"`
	ThirdCatName  string `json:"third_cat_name"`
}

// AccountGetCategoryList 获取商家类目列表
//
// 获取已申请成功的类类目列表
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/account/category_list.html
func (s *ShopComponentShopService) AccountGetCategoryList(ctx context.Context) ([]ShopCategory, error) {
	u, err := s.client.apiURL(ctx, "shop/account/get_category_list", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return nil, err
	}
	var data struct {
		Data []ShopCategory `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}

// ShopBrandListItem 品牌列表项
type ShopBrandListItem struct {
	BrandID      int
	BrandWording string
}

// AccountGetBrandList 获取商家品牌列表
//
// 获取已申请成功的品牌列表
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/account/brand_list.html
func (s *ShopComponentShopService) AccountGetBrandList(ctx context.Context) ([]ShopBrandListItem, error) {
	u, err := s.client.apiURL(ctx, "shop/account/get_brand_list", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return nil, err
	}
	var data struct {
		Data []ShopBrandListItem `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}

type ShopAccountInfo struct {
	ServiceAgentPath  string
	ServiceAgentPhone string
}

// AccountGetInfo 获取商家信息
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/account/get_info.html
func (s *ShopComponentShopService) AccountGetInfo(ctx context.Context) (*ShopAccountInfo, error) {
	u, err := s.client.apiURL(ctx, "shop/account/get_info", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return nil, err
	}
	var data struct {
		Data *ShopAccountInfo `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}
