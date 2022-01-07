package glide

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type impersonating struct {
	sub                string
	scopes             []string
	accessToken        string
	accessTokenExpires time.Time
}

type clientOptions struct {
	protocol string
	host     string
	url      string
	basePath string
	audience string
}

type clientOption func(options *clientOptions)

func WithProtocol(protocol string) clientOption {
	return func(o *clientOptions) {
		o.protocol = strings.ToLower(protocol)
	}
}

func WithURL(URL string) clientOption {
	return func(o *clientOptions) {
		o.url = URL
	}
}

func WithHost(host string) clientOption {
	return func(o *clientOptions) {
		o.host = host
	}
}

func WithBasePath(basePath string) clientOption {
	return func(o *clientOptions) {
		o.basePath = basePath
	}
}

func WithAudience(audience string) clientOption {
	return func(o *clientOptions) {
		o.audience = audience
	}
}

func WithQueryParam(qParam string, value string) requestOption {
	return withQueryParam(qParam, value)
}

func WithQueryParamList(qParam string, values ...string) requestOption {
	return withQueryParamList(qParam, values)
}

func WithExpand(paths ...string) requestOption {
	return withQueryParamList("expand", paths)
}

func GetExpandFields(fieldIds ...string) string {
	if len(fieldIds) > 0 {
		return fmt.Sprintf("fields[%s]", strings.Join(fieldIds, ","))
	}
	return "fields"
}

func WithUpdatedAfter(ts int) requestOption {
	return withQueryParam("updated_after", strconv.Itoa(ts))
}

type Client interface {
	get(res Response, authRequired bool, path string, opts ...requestOption) error
	post(res Response, authRequired bool, path string, payload Request, opts ...requestOption) error
	StartImpersonating(sub string, scopes []string) error
	IsImpersonating() bool
	ImpersonatingSub() string
	ImpersonatingScopes() []string
	StopImpersonating()
	// DO NOT remove these comments since they serve as anchors for code autogeneration
	/* Autogenerated-root-resource-interface-defs begins */
	Documents() DocumentsResource
	Listings() ListingsResource
	Notifications() NotificationsResource
	Transactions() TransactionsResource
	Users() UsersResource
	/* Autogenerated-root-resource-interface-defs ends */
}

type client struct {
	clientKey     string
	key           Key
	options       *clientOptions
	impersonating *impersonating

	// DO NOT remove these comments since they serve as anchors for code autogeneration
	/* Autogenerated-root-resource-defs begins */
	documents     DocumentsResource
	listings      ListingsResource
	notifications NotificationsResource
	transactions  TransactionsResource
	users         UsersResource
	/* Autogenerated-root-resource-defs ends */
}

const JWT_EXPIRES = 60

func GetClient(clientKey string, key Key, opts ...clientOption) Client {
	defaultOptions := clientOptions{
		protocol: "https",
		host:     "api.glide.com",
		url:      "",
		basePath: "",
		audience: "",
	}
	c := client{
		clientKey: clientKey,
		key:       key,
		options:   &defaultOptions,
	}
	for _, option := range opts {
		option(c.options)
	}
	if c.options.url == "" {
		c.options.url = c.options.host
	}
	if c.options.audience == "" {
		c.options.audience = c.options.host
	}
	// DO NOT remove these comments since they serve as anchors for code autogeneration
	/* Autogenerated-root-resource-init begins */
	c.documents = getDocumentsResource(&c)
	c.listings = getListingsResource(&c)
	c.notifications = getNotificationsResource(&c)
	c.transactions = getTransactionsResource(&c)
	c.users = getUsersResource(&c)
	/* Autogenerated-root-resource-init ends */
	return &c
}

// DO NOT remove these comments since they serve as anchors for code autogeneration
/* Autogenerated-root-resource-getters begins */
func (c client) Documents() DocumentsResource {
	return c.documents
}
func (c client) Listings() ListingsResource {
	return c.listings
}
func (c client) Notifications() NotificationsResource {
	return c.notifications
}
func (c client) Transactions() TransactionsResource {
	return c.transactions
}
func (c client) Users() UsersResource {
	return c.users
}

/* Autogenerated-root-resource-getters ends */

func (c client) getUrl(path string) string {
	return fmt.Sprintf("%s://%s%s%s", strings.ToLower(c.options.protocol),
		strings.ToLower(c.options.url),
		strings.ToLower(c.options.basePath),
		strings.ToLower(path),
	)
}

func (c client) request(authRequired bool, doRequest func(extraOpts ...requestOption) error) error {
	if authRequired && c.IsImpersonating() {
		c.refreshAccessToken()
	}

	extraOpts := []requestOption{
		withRequestHost(c.options.host),
	}

	if c.IsImpersonating() {
		extraOpts = append(extraOpts, withHeader("Authorization", fmt.Sprintf("Bearer %s", c.impersonating.accessToken)))
	}

	return doRequest(extraOpts...)
}

func (c client) get(res Response, authRequired bool, path string, opts ...requestOption) error {
	return c.request(authRequired, func(extraOpts ...requestOption) error {
		return get(res, c.getUrl(path), append(opts, extraOpts...)...)
	})
}

func (c client) post(res Response, authRequired bool, path string, payload Request, opts ...requestOption) error {
	return c.request(authRequired, func(extraOpts ...requestOption) error {
		return post(res, c.getUrl(path), payload, append(opts, extraOpts...)...)
	})
}

func (c client) getJwt(sub string, scopes []string) (string, error) {
	return getJwt(c.key, c.clientKey, sub, c.options.audience, scopes, JWT_EXPIRES)
}

func (c client) getAccessToken(sub string, scopes []string) (*jwtOauthAccessToken, error) {
	assertionsJwt, err := c.getJwt(sub, scopes)
	if err != nil {
		return nil, &ApiError{
			Params: map[string]interface{}{
				"desc": "Error issuing assertions JWT",
			},
			Err: err,
		}
	}

	data := jwtOauthAccessToken{}
	if err = c.post(&data, false, "/oauth/token", jwtOauth{
		GrantType: "JWT",
		Assertion: assertionsJwt,
	}); err != nil {
		return nil, err
	}

	if data.AccessToken == "" {
		if data.RequestScopesUrl != "" {
			err = &ApiError{
				Description: "Need user consent to access scopes",
				Params: map[string]interface{}{
					"missing_scopes":     data.MissingScopes,
					"request_scopes_url": data.RequestScopesUrl,
				},
			}
		} else {
			err = &ApiError{
				Description: "Unknown error",
				Params: map[string]interface{}{
					"data": data,
				},
			}
		}
		return nil, err
	}

	return &data, nil
}

func (c *client) StartImpersonating(sub string, scopes []string) error {
	data, err := c.getAccessToken(sub, scopes)
	if err != nil {
		return err
	}
	c.impersonating = &impersonating{
		sub:                sub,
		scopes:             scopes,
		accessToken:        data.AccessToken,
		accessTokenExpires: time.Now().Add(time.Second * time.Duration(data.ExpiresIn)),
	}

	return nil
}

func (c client) IsImpersonating() bool {
	return c.impersonating != nil
}

func (c client) ImpersonatingSub() string {
	if c.IsImpersonating() {
		return c.impersonating.sub
	}
	return ""
}

func (c client) ImpersonatingScopes() []string {
	if c.IsImpersonating() {
		return c.impersonating.scopes
	}
	return []string{}
}

func (c *client) StopImpersonating() {
	c.impersonating = nil
}

func (c client) isTokenExpired() bool {
	return time.Now().After(c.impersonating.accessTokenExpires)
}

func (c client) isTokenAboutToExpire() bool {
	// Consider an access token "about to expire" if it's set to expire in less than 15m
	return !c.isTokenExpired() && time.Now().After(c.impersonating.accessTokenExpires.Add(time.Minute*time.Duration(-15)))
}

func (c client) shouldRefreshToken() bool {
	return c.IsImpersonating() && (c.isTokenAboutToExpire() || c.isTokenExpired())
}

func (c *client) refreshAccessToken() error {
	if c.shouldRefreshToken() {
		data, err := c.getAccessToken(c.impersonating.sub, c.impersonating.scopes)

		if err != nil {
			// fail silently if token is not yet expired
			if c.isTokenExpired() {
				return err
			}
		} else {
			c.impersonating.accessToken = data.AccessToken
			c.impersonating.accessTokenExpires = time.Now().Add(time.Second * time.Duration(data.ExpiresIn))
		}
	}

	return nil
}
