package miniprogram

import (
	"context"
	"fmt"

	"github.com/dawndiy/go-wechat/pkg/upload"
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
	BrandID      int    `json:"brand_id"`
	BrandWording string `json:"brand_wording"`
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
	ServiceAgentPath        string                                  `json:"service_agent_path"`
	ServiceAgentPhone       string                                  `json:"service_agent_phone"`
	ServiceAgentType        []int                                   `json:"service_agent_type"`
	DefaultReceivingAddress *ShopAccountInfoDefaultReceivingAddress `json:"default_receiving_address"`
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

// ImageUpload 上传图片
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/public/upload_img.html
func (s *ShopComponentShopService) ImageUpload(ctx context.Context, up upload.Uploader, respType int) (string, error) {
	u, err := s.client.apiURL(ctx, "shop/img/upload", nil)
	if err != nil {
		return "", err
	}
	fields := map[string]string{
		"resp_type":   fmt.Sprint(respType),
		"upload_type": "0", //0:图片流，1:图片url
	}
	req, err := s.client.NewUploadRequest(ctx, u.String(), "media", up.Name(), up.Reader(), fields)
	if err != nil {
		return "", err
	}
	var data struct {
		ImgInfo struct {
			TempImgURL string `json:"temp_img_url"`
			MediaID    string `json:"media_id"`
		} `json:"img_info"`
	}
	_, err = s.client.Do(req, &data)
	if err != nil {
		return "", err
	}
	if respType == 0 {
		return data.ImgInfo.MediaID, nil
	}

	return data.ImgInfo.TempImgURL, nil
}

// ShopAuditBrandRequest 自定义交易组件上传品牌请求
type ShopAuditBrandRequest struct {
	// 营业执照或组织机构代码证，图片url/media_id
	License []string `json:"license"`
	// 品牌信息
	BrandInfo ShopAuditBrandInfo `json:"brand_info"`
	// 商品使用场景,1:视频号，3:订单中心
	SceneGroupList []int `json:"scene_group_list"`
}

// ShopAuditBrandInfo 自定义交易组件上传品牌信息
type ShopAuditBrandInfo struct {
	// 认证审核类型
	BrandAuditType int `json:"brand_audit_type"`
	// 商标分类
	TrademarkType string `json:"trademark_type"`
	// 选择品牌经营类型
	BrandManagementType int `json:"brand_management_type"`
	// 商品产地是否进口
	CommodityOriginType int `json:"commodity_origin_type"`
	// 商标/品牌词
	BrandWording string `json:"brand_wording"`
	// 销售授权书（如商持人为自然人，还需提供有其签名的身份证正反面扫描件)，
	// 图片url/media_id
	SaleAuthorization []string `json:"sale_authorization,omitempty"`
	// 商标注册证书，图片url/media_id
	TrademarkRegistrationCertificate []string `json:"trademark_registration_certificate,omitempty"`
	// 商标变更证明，图片url/media_id
	TrademarkChangeCertificate []string `json:"trademark_change_certificate,omitempty"`
	// 商标注册人姓名
	TrademarkRegistrant string `json:"trademark_registrant,omitempty"`
	// 商标注册号/申请号
	TrademarkRegistrantNu string `json:"trademark_registrant_nu"`
	// 商标有效期，yyyy-MM-dd HH:mm:ss
	TrademarkAuthorizationPeriod string `json:"trademark_authorization_period,omitempty"`
	// 商标注册申请受理通知书，图片url/media_id
	TrademarkRegistrationApplication []string `json:"trademark_registration_application,omitempty"`
	// 商标申请人姓名
	TrademarkApplicant string `json:"trademark_applicant,omitempty"`
	// 商标申请时间, yyyy-MM-dd HH:mm:ss
	TrademarkApplicationTime string `json:"trademark_application_time,omitempty"`
	// 中华人民共和国海关进口货物报关单，图片url/media_id
	ImportedGoodsForm []string `json:"imported_goods_form,omitempty"`
}

// AuditBrand 品牌审核
//
// 请求成功后将会创建一个审核单，单号将在回包中给出。
// 审核完成后会进行回调，告知审核结果，请接入品牌审核回调，如果审核成功，则在回调中给出brand_id。
// 使用到的图片的地方，可以使用url或media_id(通过上传图片接口换取)。
// 请确认图片url可以正常打开，图片大小在2MB以下，图片格式为jpg, jpeg, png，如图片不能正常显示，会导致审核驳回。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/audit/audit_brand.html
func (s *ShopComponentShopService) AuditBrand(ctx context.Context, audit ShopAuditBrandRequest) (string, error) {
	u, err := s.client.apiURL(ctx, "shop/audit/audit_brand", nil)
	if err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"audit_req": audit,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}
	var data struct {
		AuditID string `json:"audit_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.AuditID, err
}

// ShopAuditCategoryRequest 自定义交易组件类目审核
type ShopAuditCategoryRequest struct {
	// 营业执照或组织机构代码证，图片url
	License []string `json:"license"`
	// 类目信息
	CategoryInfo ShopAuditCategoryInfo `json:"category_info"`
}

// ShopAuditCategoryInfo 自定义交易组件类目审核信息
type ShopAuditCategoryInfo struct {
	// 一级类目
	Level1 int `json:"level1"`
	// 二级类目
	Level2 int `json:"level2"`
	// 三级类目
	Level3 int `json:"level3"`
	// 资质材料，图片url
	Certificate []string `json:"certificate"`
}

// AuditCategory 类目审核
//
// 请求成功后将会创建一个审核单，单号将在回包中给出
// 审核完成后会进行回调，告知审核结果
// 这个上传类目资质的接口，如果上传的类目是已经审核通过的，该接口会返回错误码 1050003
// 使用到的图片的地方，可以使用url或media_id(通过上传图片接口换取)。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/audit/audit_category.html
func (s *ShopComponentShopService) AuditCategory(ctx context.Context, audit ShopAuditCategoryRequest) (string, error) {
	u, err := s.client.apiURL(ctx, "shop/audit/audit_category", nil)
	if err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"audit_req": audit,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}
	var data struct {
		AuditID string `json:"audit_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.AuditID, err
}

// ShopAuditResult 查询品牌和类目的审核结果
type ShopAuditResult struct {
	// 审核状态, 0：审核中，1：审核成功，9：审核拒绝
	Status int `json:"status"`
	// 如果审核拒绝，返回拒绝原因
	RejectReason string `json:"reject_reason"`
	// 如果是品牌审核，返回brand_id
	BrandID int `json:"brand_id"`
}

// AuditResult 获取审核结果
//
// 根据审核id，查询品牌和类目的审核结果。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/audit/audit_result.html
func (s *ShopComponentShopService) AuditResult(ctx context.Context, auditID string) (*ShopAuditResult, error) {
	u, err := s.client.apiURL(ctx, "shop/audit/result", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"audit_id": auditID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data *ShopAuditResult
	}
	_, err = s.client.Do(req, &result)
	return result.Data, err
}

// ShopMiniappCertificate 自定义交易组件小程序资质
type ShopMiniappCertificate struct {
	BrandInfo struct {
		BrandWording                     string   `json:"brand_wording"`
		SaleAuthorization                []string `json:"sale_authorization"`
		TrademarkRegistrationCertificate []string `json:"trademark_registration_certificate"`
	} `json:"brand_info"`
	CategoryInfoList []struct {
		FirstCategoryID    int      `json:"first_category_id"`
		FirstCategoryName  int      `json:"first_category_name"`
		SecondCategoryID   int      `json:"second_category_id"`
		SecondCategoryName int      `json:"second_category_name"`
		CertificateURL     []string `json:"certificate_url"`
	} `json:"category_info_list"`
}

// AuditGetMiniappCertificate 获取小程序资质
//
// 获取曾经提交的小程序审核资质
// 请求类目会返回多次的请求记录，请求品牌只会返回最后一次的提交记录
// 图片经过转链，请使用高版本chrome浏览器打开
// 如果曾经没有提交，没有储存历史文件，或是获取失败，接口会返回1050006
// 注：该接口返回的是曾经在小程序方提交过的审核，非组件的入驻审核！
//
// https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/audit/get_miniapp_certificate.html
func (s *ShopComponentShopService) AuditGetMiniappCertificate(ctx context.Context, reqType int) (*ShopMiniappCertificate, error) {
	u, err := s.client.apiURL(ctx, "shop/audit/get_miniapp_certificate", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"req_type": reqType,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var data ShopMiniappCertificate
	_, err = s.client.Do(req, &data)
	return &data, err
}

// ShopAccountUpdateInfoRequest 商家信息更新请求
type ShopAccountUpdateInfoRequest struct {
	// 小程序path 可选
	ServiceAgentPath string `json:"service_agent_path,omitempty"`
	// 客服联系方式 可选
	ServiceAgentPhone string `json:"service_agent_phone,omitempty"`
	// 客服类型，支持多个，0: 小程序客服，1: 自定义客服path 2: 联系电话
	ServiceAgentType []int `json:"service_agent_type"`
	// 默认退货地址
	DefaultReceivingAddress *ShopAccountInfoDefaultReceivingAddress `json:"default_receiving_address,omitempty"`
}

type ShopAccountInfoDefaultReceivingAddress struct {
	// 收件人姓名 必填
	ReceiverName string `json:"receiver_name"`
	// 详细收货地址信息 必填
	DetailedAddress string `json:"detailed_address"`
	// 收件人手机号码 必填
	TelNumber string `json:"tel_number"`
	// 国家 可选
	Country string `json:"country,omitempty"`
	// 省份 可选
	Province string `json:"province,omitempty"`
	// 城市 可选
	City string `json:"city,omitempty"`
	// 乡镇 可选
	Town string `json:"town,omitempty"`
}

// AccountUpdateInfo 更新商家信息
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/account/update_info.html
func (s *ShopComponentShopService) AccountUpdateInfo(ctx context.Context, r ShopAccountUpdateInfoRequest) error {
	u, err := s.client.apiURL(ctx, "shop/account/update_info", nil)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), r)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
