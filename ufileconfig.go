package ufile_gosdk

//北京
//外网：www.cn-bj.ufileos.com
//B机房内网：www.ufile.cn-north-02.ucloud.cn
//C机房内网: www.ufile.cn-north-03.ucloud.cn
//D机房内网: www.ufile.cn-north-04.ucloud.cn
//上海二
//外网: www.cn-sh2.ufileos.com
//内网: www.internal-cn-sh2-01.ufileos.com
//香港
//外网：www.hk.ufileos.com
//内网：www.internal-hk-01.ufileos.com
//广东
//外网：www.cn-gd.ufileos.com
//内网：www.internal-cn-gd-02.ufileos.com
//美国
//外网：www.us-ca.ufileos.com
//内网：www.internal-us-ca-01.ufileos.com

const (
	CN_BEIJING = "www.cn-bj.ufileos.com"
	CN_BEIJING_INTERNAL_B = "www.ufile.cn-north-02.ucloud.cn"
	CN_BEIJING_INTERNAL_C = "www.ufile.cn-north-03.ucloud.cn"
	CN_BEIJING_INTERNAL_D = "www.ufile.cn-north-04.ucloud.cn"

	CN_SHANGHAI = "www.cn-sh2.ufileos.com"
	CN_SHANGHAI_INTERNAL = "www.internal-cn-sh2-01.ufileos.com"

	CN_HK = "www.hk.ufileos.com"
	CN_HK_INTERNAL = "www.internal-hk-01.ufileos.com"

	CN_GUANGDONG = "www.cn-gd.ufileos.com"
	CN_GUANGDONG_INTERNAL = "www.internal-cn-gd-02.ufileos.com"

	US_CA = "www.us-ca.ufileos.com"
	US_CA_INTERNAL = "www.internal-us-ca-01.ufileos.com"
)

var Region map[string]string

func init() {
	Region = make(map[string]string)
	Region[CN_BEIJING] = "cn-bj2"
	Region[CN_BEIJING_INTERNAL_B] = "cn-bj2"
	Region[CN_BEIJING_INTERNAL_C] = "cn-bj2"
	Region[CN_BEIJING_INTERNAL_D] = "cn-bj2"

	Region[CN_SHANGHAI] = "cn-sh2"
	Region[CN_SHANGHAI_INTERNAL] = "cn-sh2"

	Region[CN_HK] = "hk"
	Region[CN_HK_INTERNAL] = "hk"

	Region[CN_GUANGDONG] = "cn-gd"
	Region[CN_GUANGDONG_INTERNAL] = "cn-gd"

	Region[US_CA] = "us-ca"
	Region[US_CA_INTERNAL] = "us-ca"
}