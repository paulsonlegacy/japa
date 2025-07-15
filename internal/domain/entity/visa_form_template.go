package entity


type VisaFormTemplate struct {
    ID             uint    	`gorm:"type:int;primaryKey"`
    Country        string   
    VisaType       string
    DownloadURL    string    // Where user can get the form
    Instructions   string    // Extra notes on how to fill
}