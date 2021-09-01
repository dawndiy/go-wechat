package miniprogram

import (
	"context"
)

// ShopSPUAddInfo 交易组件商品
type ShopSPUAddInfo struct {
	// 商家自定义商品ID 必填
	OutProductID string `json:"out_product_id"`
	// 标题 必填
	Title string `json:"title"`
	// 绑定的小程序商品路径 必填
	Path string `json:"path"`
	// 主图,多张,列表 必填
	HeadImg []string `json:"head_img"`
	// 商品资质图片 选填
	QualificationPics []string `json:"qualification_pics,omitempty"`
	// 商品详情描述 选填
	DescInfo *ShopSPUDescInfo `json:"desc_info,omitempty"`
	// 第三级类目ID 必填
	ThirdCatID int `json:"third_cat_id"`
	// 品牌id 必填
	BrandID int `json:"brand_id"`
	// 预留字段，用于版本控制 选填
	InfoVersion string `json:"info_version"`
	// sku数组
	SKUS []ShopSKU `json:"skus"`
}

// ShopSPUDescInfo 商品描述信息
type ShopSPUDescInfo struct {
	Desc string   `json:"desc"`
	Imgs []string `json:"imgs"`
}

// ShopSKU 商品库存
type ShopSKU struct {
	// 商家自定义商品ID 必填
	OutProductID string `json:"out_product_id"`
	// 商家自定义skuID 必填
	OutSKUID string `json:"out_sku_id"`
	// sku小图 必填
	ThumbImg string `json:"thumb_img"`
	// 售卖价格,以分为单位 必填
	SalePrice int64 `json:"sale_price"`
	// 市场价格,以分为单位 必填
	MarketPrice int64 `json:"market_price"`
	// 库存 必填
	StockNum int64 `json:"stock_num"`
	// 条形码 选填
	Barcode string `json:"barcode"`
	// 商品编码 选填
	SKUCode string `json:"sku_code"`
	// 销售属性 选填
	SKUAttrs []ShopSPUAttr `json:"sku_attrs"`
}

// ShopSPUAttr 销售属性
type ShopSPUAttr struct {
	// 销售属性key（自定义）
	AttrKey string `json:"attr_key"`
	// 销售属性value（自定义）
	AttrValue string `json:"attr_value"`
}

type ShopSPUAddResult struct {
	// 交易组件平台内部商品ID
	ProductID    int64  `json:"product_id"`
	OutProductID string `json:"out_product_id"`
	CreateTime   string `json:"create_time"`
	SKUS         []struct {
		SKUID    int64  `json:"sku_id"`
		OutSKUID string `json:"out_sku_id"`
	} `json:"skus"`
}

// SPUAdd 添加商品
//
// 新增成功后会直接提交审核，可通过商品审核回调，或者通过get接口的edit_status查看是否通过审核。
// 商品有2份数据，草稿和线上数据。
// 调用接口新增或修改商品数据后，影响的只是草稿数据，审核通过后草稿数据才会覆盖线上数据，正式生效。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/add_spu.html
func (s *ShopComponentShopService) SPUAdd(ctx context.Context, spu ShopSPUAddInfo) (*ShopSPUAddResult, error) {
	u, err := s.client.apiURL(ctx, "shop/spu/add", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), spu)
	if err != nil {
		return nil, err
	}
	var data struct {
		Data *ShopSPUAddResult
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}

// SPUDel 删除商品
//
// 从初始值/上架/若干下架状态转换成逻辑删除（删除后不可恢复）
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/del_spu.html
func (s *ShopComponentShopService) SPUDel(ctx context.Context, productID int64, outProductID string) error {
	u, err := s.client.apiURL(ctx, "shop/spu/del", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"product_id":     productID,
		"out_product_id": outProductID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// SPUDelAudit 撤回商品审核
//
// 对于审核中（edit_status=2）的商品无法重复提交，需要调用此接口，使商品流转进入未审核的状态（edit_status=1）,即可重新提交商品。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/del_spu_audit.html
func (s *ShopComponentShopService) SPUDelAudit(ctx context.Context, productID int64, outProductID string) error {
	u, err := s.client.apiURL(ctx, "shop/spu/del_audit", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"product_id":     productID,
		"out_product_id": outProductID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// ShopSPU 交易组件商品
type ShopSPU struct {
	// 交易组件平台内部商品ID
	ProductID int64 `json:"product_id"`
	// 商家自定义商品ID 必填
	OutProductID string `json:"out_product_id"`
	// 标题 必填
	Title string `json:"title"`
	// 绑定的小程序商品路径 必填
	Path string `json:"path"`
	// 主图,多张,列表 必填
	HeadImg []string `json:"head_img"`
	// 商品详情描述 选填
	DescInfo *ShopSPUDescInfo `json:"desc_info,omitempty"`
	// 商品审核信息
	AuditInfo *ShopSPUAuditInfo `json:"audit_info,omitempty"`
	// 商品线上状态
	Status int `json:"status"`
	// 商品草稿状态
	EditStatus int `json:"edit_status"`
	// 第三级类目ID
	ThirdCatID int `json:"third_cat_id"`
	// 品牌id 必填
	BrandID int `json:"brand_id"`
	// 创建时间
	CreateTime string `json:"create_time"`
	// 更新时间
	UpdateTime string `json:"update_time"`
	// 预留字段，用于版本控制 选填
	InfoVersion string `json:"info_version"`
	// sku数组
	SKUS []ShopSKU `json:"skus"`
}

// ShopSPUAuditInfo 商品审核信息
type ShopSPUAuditInfo struct {
	SubmitTime   string `json:"submit_time"`
	AuditTime    string `json:"audit_time"`
	RegectReason string `json:"regect_reason"`
	AuditID      string `json:"audit_id"`
}

// SPUGet 获取商品
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/get_spu.html
func (s *ShopComponentShopService) SPUGet(ctx context.Context, productID int64, outProductID string, needEditSPU int) (*ShopSPU, error) {
	u, err := s.client.apiURL(ctx, "shop/spu/get", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"product_id":     productID,
		"out_product_id": outProductID,
		"need_edit_spu":  needEditSPU,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var data struct {
		SPU *ShopSPU
	}
	_, err = s.client.Do(req, &data)
	return data.SPU, err
}

// ShopSPUListQuery 交易组件商品列表查询
type ShopSPUListQuery struct {
	// 商品状态 选填，不填时获取所有状态商品
	Status int `json:"status,omitempty"`
	// 开始创建时间 选填，与end_create_time成对
	StartCreateTime string `json:"start_create_time,omitempty"`
	// 结束创建时间 选填，与start_create_time成对
	EndCreateTime string `json:"end_create_time,omitempty"`
	// 开始更新时间 选填，与end_update_time成对
	StartUpdateTime string `json:"start_update_time,omitempty"`
	// 结束更新时间 选填，与start_update_time成对
	EndUpdateTime string `json:"end_update_time,omitempty"`
	// 页号
	Page int64 `json:"page"`
	// 页面大小 不超过100
	PageSize int64 `json:"page_size"`
}

// SPUGetList 获取商品列表
//
// 时间范围 create_time 和 update_time 同时存在时，以 create_time 的范围为准
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/get_spu_list.html
func (s *ShopComponentShopService) SPUGetList(ctx context.Context, productID int64, outProductID string, needEditSPU int) (*ShopSPU, error) {
	u, err := s.client.apiURL(ctx, "shop/spu/get_list", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"product_id":     productID,
		"out_product_id": outProductID,
		"need_edit_spu":  needEditSPU,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var data struct {
		SPU *ShopSPU
	}
	_, err = s.client.Do(req, &data)
	return data.SPU, err
}

// ShopSPUUpdateInfo 交易组件商品更新信息
type ShopSPUUpdateInfo struct {
	ShopSPUAddInfo `json:",inline"`
	ProductID      int64 `json:"product_id"`
}

// ShopSPUUpdateResult 交易组件商品更新结果
type ShopSPUUpdateResult struct {
	// 交易组件平台内部商品ID
	ProductID    int64  `json:"product_id"`
	OutProductID string `json:"out_product_id"`
	UpdateTime   string `json:"create_time"`
	SKUS         []struct {
		SKUID    int64  `json:"sku_id"`
		OutSKUID string `json:"out_sku_id"`
	} `json:"skus"`
}

// SPUUpdate 更新商品
//
// 注意：更新成功后会更新到草稿数据并直接提交审核，审核完成后有回调，也可通过get接口的edit_status查看是否通过审核。
// 商品有两份数据，草稿和线上数据。
// 调用接口新增或修改商品数据后，影响的只是草稿数据，审核通过草稿数据才会覆盖线上数据正式生效。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/update_spu.html
func (s *ShopComponentShopService) SPUUpdate(ctx context.Context, spu ShopSPUAddInfo) (*ShopSPUUpdateResult, error) {
	u, err := s.client.apiURL(ctx, "shop/spu/update", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), spu)
	if err != nil {
		return nil, err
	}
	var data struct {
		Data *ShopSPUUpdateResult
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}

// SPUListing 上架商品
//
// 如果该商品处于自主下架状态，调用此接口可把直接把商品重新上架 该接口不影响已经在审核流程的草稿数据
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/listing_spu.html
func (s *ShopComponentShopService) SPUListing(ctx context.Context, productID int64, outProductID string) error {
	u, err := s.client.apiURL(ctx, "shop/spu/listing", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"product_id":     productID,
		"out_product_id": outProductID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// SPUDelisting 下架商品
//
// 从初始值/上架状态转换成自主下架状态
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/delisting_spu.html
func (s *ShopComponentShopService) SPUDelisting(ctx context.Context, productID int64, outProductID string) error {
	u, err := s.client.apiURL(ctx, "shop/spu/delisting", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"product_id":     productID,
		"out_product_id": outProductID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// ShopSPUUpdateWithoutAuditInfo 免审更新商品字段
type ShopSPUUpdateWithoutAuditInfo struct {
	OutProductID string    `json:"out_product_id"`
	ProductID    int64     `json:"product_id"`
	Path         string    `json:"path"`
	SKUS         []ShopSKU `json:"skus,omitempty"`
}

// SPUUpdateWithoutAudit 更新商品
//
// 注意：该免审接口只能更新部分商品字段，影响草稿数据和线上数据，
// 且请求包中的sku必须已经存在于原本的在线数据中（比如out_sku_id="123"如果不在原本的线上数据的sku列表中，将返回错误1000004）
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/SPU/update_spu_without_audit.html
func (s *ShopComponentShopService) SPUUpdateWithoutAudit(ctx context.Context, spu ShopSPUUpdateWithoutAuditInfo) (*ShopSPUUpdateResult, error) {
	u, err := s.client.apiURL(ctx, "shop/spu/update_without_audit", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), spu)
	if err != nil {
		return nil, err
	}
	var data struct {
		Data *ShopSPUUpdateResult
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}
