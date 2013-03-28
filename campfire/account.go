package campfire

type Account struct {
    *Client

    Id int `json:"id"`
    OwnerId string `json:"owner_id"`
    Subdomain string `json:"subdomain"`
    Name string `json:"name"`
    TimeZone string `json:"time_zone"`
    Storage int `json:"storage"` 
    Plan string `json:"plan"`
}
