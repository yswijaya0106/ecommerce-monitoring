package customer

// Customer mirrors the `customers` table in the Northwind schema, minus the
// `attachments` BLOB column which is not useful over a JSON API.
type Customer struct {
	ID            int64   `json:"id"`
	Company       *string `json:"company,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	FirstName     *string `json:"first_name,omitempty"`
	EmailAddress  *string `json:"email_address,omitempty"`
	JobTitle      *string `json:"job_title,omitempty"`
	BusinessPhone *string `json:"business_phone,omitempty"`
	HomePhone     *string `json:"home_phone,omitempty"`
	MobilePhone   *string `json:"mobile_phone,omitempty"`
	FaxNumber     *string `json:"fax_number,omitempty"`
	Address       *string `json:"address,omitempty"`
	City          *string `json:"city,omitempty"`
	StateProvince *string `json:"state_province,omitempty"`
	ZipPostalCode *string `json:"zip_postal_code,omitempty"`
	CountryRegion *string `json:"country_region,omitempty"`
	WebPage       *string `json:"web_page,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}
