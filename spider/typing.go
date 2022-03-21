package spider

const (
	chinaTotalUri        string = "https://view.inews.qq.com/g2/getOnsInfo?name=disease_h5"
	provinceCityTotalUri string = "https://api.inews.qq.com/newsqa/v1/query/inner/publish/modules/list?modules=statisGradeCityDetail,diseaseh5Shelf"
)

type ChinaTotal struct {
	LastUpdateTime string `json:"lastUpdateTime"`
	Confirm        int64  `json:"confirm"`
	Dead           int64  `json:"dead"`
	NowConfirm     int64  `json:"nowConfirm"`
	ImportedCase   int64  `json:"importedCase"`
	NoInfect       int64  `json:"noInfect"`
	LocalConfirm   int64  `json:"localConfirm"`
}

type ProvinceCityTotal struct {
	NowConfirm int64 `json:"nowConfirm"`
	Wzz        int64 `json:"wzz"`
	Heal       int64 `json:"heal"`
	Confirm    int64 `json:"confirm"`
	Dead       int64 `json:"dead"`
	Exists     bool
}
