package db

type AuthenticationProvider string

const (
	AuthProviderPassword = "password"
	AuthProviderGoogle   = "google"
)

type MembershipStatus string

func (s MembershipStatus) String() string {
	return string(s)
}

const (
	MemberStatusActive  MembershipStatus = "active"
	MemberStatusInctive MembershipStatus = "inactive"
)

type InvitationStatus string

func (s InvitationStatus) String() string {
	return string(s)
}

const (
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusCanceled InvitationStatus = "canceled"
	InvitationStatusPending  InvitationStatus = "pending"
)

type ChannelProvider string

func (p ChannelProvider) String() string {
	return string(p)
}

const (
	ChannelProviderShopify     = "shopify"
	ChannelProviderWooCommerce = "woocommerce"
)

type ChannelStatus string

func (s ChannelStatus) String() string {
	return string(s)
}

const (
	ChannelStatusActive   ChannelStatus = "active"
	ChannelStatusInactive ChannelStatus = "inactive"
)
