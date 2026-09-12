package db

type AuthenticationProvider string

const (
	AuthProviderPassword = "password"
	AuthProviderGoogle   = "google"
)

type MembershipRole string

func (r MembershipRole) String() string {
	return string(r)
}

const (
	RoleOwner  MembershipRole = "owner"
	RoleAdmin  MembershipRole = "admin"
	RoleMember MembershipRole = "member"
)

type MembershipPermission string

func (p MembershipPermission) String() string {
	return string(p)
}

const (
	PermMembershipRead   MembershipPermission = "membership:read"
	PermMembershipCreate MembershipPermission = "membership:create"
	PermMembershipUpdate MembershipPermission = "membership:update"
	PermMembershipDelete MembershipPermission = "membership:delete"
	PermChannelRead      MembershipPermission = "channel:read"
	PermChannelCreate    MembershipPermission = "channel:create"
	PermChannelUpdate    MembershipPermission = "channel:update"
	PermChannelDelete    MembershipPermission = "channel:delete"
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

// type ChannelCredentials struct

// Shopify
// StoreURL
// AccessToken
// APIKey
// SecretKey

// Woo
// StoreURL
// ConsumerKey
// ConsumerSecret

type ChannelStatus string

func (s ChannelStatus) String() string {
	return string(s)
}

const (
	ChannelStatusActive   ChannelStatus = "active"
	ChannelStatusInactive ChannelStatus = "inactive"
)
