package data

import (
	"database/sql"
	"encoding/json"
	"image"
	"log"
	"math/big"
	"strings"

	"github.com/AlexNa-Holdings/savva-reports/assets"
	"github.com/AlexNa-Holdings/savva-reports/cmn"
)

type UserProfile struct {
	About         string `json:"about"`
	DisplayName   string `json:"display_name"`
	SponsorValues []int  `json:"sponsor_values"`
	Name          string `json:"name"`
}

type User struct {
	Address    string
	Name       string
	AvatarCid  string
	AvatarData []byte // In-memory avatar image data
	AvatarImg  image.Image
	Staked     big.Int
	Profiles   map[string]UserProfile // domaion -> profile
}

var UserCache = cmn.NewCache[string, *User](100)

func GetUser(address string) (*User, error) {
	// Check if the user is already cached
	if o, found := UserCache.Get(address); found {
		return o, nil
	}

	user := User{
		Address:    address,
		Staked:     *big.NewInt(0),
		Profiles:   make(map[string]UserProfile),
		AvatarData: assets.AvatarDefault,
	}

	var n_staked sql.NullString
	var n_avatar sql.NullString
	var n_name sql.NullString

	err := cmn.C.DB.QueryRow(`SELECT name,avatar,staked FROM users WHERE user_addr = $1`, address).Scan(
		&n_name, &n_avatar, &n_staked)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		log.Printf("Error querying user %s: %v", address, err)
		return nil, err
	}

	user.AvatarCid = n_avatar.String
	user.Name = n_name.String
	if n_staked.Valid {
		user.Staked.SetString(n_staked.String, 10)
	}

	// load the avatar from IPFS
	if user.AvatarCid != "" {
		user.AvatarImg, err = cmn.LoadImage(user.AvatarCid)
	}

	if user.AvatarImg == nil {
		user.AvatarImg = assets.AvatarDefaultImg
	}

	rows, err := cmn.C.DB.Query(`SELECT domain, value FROM user_params WHERE user_addr = $1 AND key = 'profile_cid'`, address)
	if err != nil {
		log.Printf("Error querying user %s: %v", address, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var domain string
		var value string
		err := rows.Scan(&domain, &value)
		if err != nil {
			log.Printf("Error scanning user %s: %v", address, err)
			return nil, err
		}

		// load the profile from IPFS
		data := cmn.C.IPFS(value)
		if data == nil {
			log.Printf("Error loading IPFS file for user %s: %s", address, value)
			return nil, err
		}
		// unmarshal the profile
		profile := UserProfile{}
		err = json.Unmarshal(data, &profile)
		if err != nil {
			log.Printf("Error unmarshalling user %s: %v", address, err)
			return nil, err
		}
		user.Profiles[domain] = profile
	}

	return &user, nil
}

func (u *User) GetDisplayName() string {

	p, ok := u.Profiles["savva.app"]
	if ok {
		if p.DisplayName != "" {
			return p.DisplayName
		}
	}

	for _, p := range u.Profiles {
		if p.DisplayName != "" {
			return p.DisplayName
		}
	}

	return ""
}

func (u *User) BestName() string {
	if u.Name != "" {
		return strings.ToUpper(u.Name) + "\u00AE"
	}

	dn := u.GetDisplayName()
	if dn != "" {
		return dn
	}

	return u.Address[0:6] + "..." + u.Address[len(u.Address)-4:]
}
