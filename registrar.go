package main

import (
	godaddy "github.com/eufelipemateus/registrar-domain/registrars/godaddy"
	name "github.com/eufelipemateus/registrar-domain/registrars/name"
)

type RegistrarAuthList struct {
	GodaddyAPIKey godaddy.APIKey
	NameAPIKey    name.APIKey
}

type Domain struct {
	DomainName      string `json:"domain"`
	RegistrarIanaid string `json:"registrar_iana_id"`

	RegistrarAuthList RegistrarAuthList
	isProduction      bool
}

func (domain Domain) Renew(period int) error {

	switch domain.RegistrarIanaid {

	case "146":
		registar := godaddy.Registrar{
			Domain:       domain.DomainName,
			Period:       period,
			ApiKey:       domain.RegistrarAuthList.GodaddyAPIKey,
			IsProduction: domain.isProduction,
		}
		

		return registar.Renew()
	case "625":
		registar := name.Registrar{
			Domain:       domain.DomainName,
			Period:       period,
			APIKey:       domain.RegistrarAuthList.NameAPIKey,
			IsProduction: domain.isProduction,
		}
		return registar.Renew()
	}

	return nil
}
