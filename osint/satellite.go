package osint

type Satellite struct {
    INTLDES      string  `json:"INTLDES"`
    NORAD_CAT_ID string  `json:"NORAD_CAT_ID"`
    OBJECT_TYPE  string  `json:"OBJECT_TYPE"`
    SATNAME      string  `json:"SATNAME"`
    COUNTRY      string  `json:"COUNTRY"`
    LAUNCH       string  `json:"LAUNCH"`
    SITE         string  `json:"SITE"`
    DECAY        *string `json:"DECAY,omitempty"`
    PERIOD       float64 `json:"PERIOD,string"`
    INCLINATION  float64 `json:"INCLINATION,string"`
    APOGEE       int     `json:"APOGEE,string"`
    PERIGEE      int     `json:"PERIGEE,string"`
    COMMENT      *string `json:"COMMENT,omitempty"`
    COMMENTCODE  *string `json:"COMMENTCODE,omitempty"`
    RCSVALUE     string  `json:"RCSVALUE"`
    RCS_SIZE     *string `json:"RCS_SIZE,omitempty"`
    FILE         string  `json:"FILE"`
    LAUNCH_YEAR  string  `json:"LAUNCH_YEAR"`
    LAUNCH_NUM   string  `json:"LAUNCH_NUM"`
    LAUNCH_PIECE string  `json:"LAUNCH_PIECE"`
    CURRENT      string  `json:"CURRENT"`
    OBJECT_NAME  string  `json:"OBJECT_NAME"`
    OBJECT_ID    string  `json:"OBJECT_ID"`
    OBJECT_NUMBER string `json:"OBJECT_NUMBER"`
}
