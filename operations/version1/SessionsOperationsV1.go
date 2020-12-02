package operations1

import (
	"context"
	"net/http"

	sessclients1 "github.com/pip-services-users/pip-clients-sessions-go/version1"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	rpcservices "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type SessionsOperationsV1 struct {
	*rpcservices.RestOperations
	defaultConfig *cconf.ConfigParams

	cookie        string
	cookieEnabled bool
	maxCookieAge  int64

	accountsClient  accclients1.IAccountsClientV1
	sessionsClient  sessclients1.ISessionsClientV1
	passwordsClient passclients1.IPasswordsClientV1
	rolesClient     roleclients1.IRolesClientV1
}

func NewSessionsOperationsV1() *SessionsOperationsV1 {
	c := SessionsOperationsV1{
		RestOperations: rpcservices.NewRestOperations(),
	}

	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"options.cookie_enabled", true,
		"options.cookie", "x-session-id",
		"options.max_cookie_age", 365*24*60*60*1000,
	)
	c.cookie = "x-session-id"
	c.cookieEnabled = true
	c.maxCookieAge = 365 * 24 * 60 * 60 * 1000

	c.DependencyResolver.Put("accounts", cref.NewDescriptor("pip-services-accounts", "client", "*", "*", "1.0"))
	c.DependencyResolver.Put("passwords", cref.NewDescriptor("pip-services-passwords", "client", "*", "*", "1.0"))
	c.DependencyResolver.Put("roles", cref.NewDescriptor("pip-services-roles", "client", "*", "*", "1.0"))
	c.DependencyResolver.Put("sessions", cref.NewDescriptor("pip-services-sessions", "client", "*", "*", "1.0"))
}

func (c *SessionsOperationsV1) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.DependencyResolver.Configure(config)

	c.cookieEnabled = config.GetAsBooleanWithDefault("options.cookie_enabled", c.cookieEnabled)
	c.cookie = config.GetAsStringWithDefault("options.cookie", c.cookie)
	c.maxCookieAge = config.GetAsLongWithDefault("options.max_cookie_age", c.maxCookieAge)
}

func (c *SessionsOperationsV1) SetReferences(references cref.IReferences) {
	c.RestOperations.SetReferences(references)

	dependency, _ := c.DependencyResolver.GetOneRequired("sessions")
	sesionsClient, ok1 := dependency.(sessclients1.ISessionsClientV1)
	if !ok1 {
		panic("SessionOperationsV1: Cant't resolv dependency 'client' to ISessionsClientV1")
	}
	c.sessionsClient = sesionsClient

	dependency, _ = c.DependencyResolver.GetOneRequired("accounts")
	acountClient, ok2 := dependency.(accclients1.IAccountsClientV1)
	if !ok2 {
		panic("SessionOperationsV1: Cant't resolv dependency 'client' to IAccountsClientV1")
	}
	c.accountsClient = acountClient

	dependency, _ = c.DependencyResolver.GetOneRequired("passwords")
	passClient, ok3 := dependency.(passclients1.IPasswordsClientV1)
	if !ok3 {
		panic("SessionOperationsV1: Cant't resolv dependency 'client' to IPasswordsClientV1")
	}
	c.passwordsClient = passClient

	dependency, _ = c.DependencyResolver.GetOneRequired("roles")
	rolesClient, ok4 := dependency.(roleclients1.IRolesClientV1)
	if !ok4 {
		panic("SessionOperationsV1: Cant't resolv dependency 'client' to IRolesClientV1")
	}
	c.rolesClient = rolesClient
}

func (c *SessionsOperationsV1) LoadSession(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	sessionId := req.Header.Get("x-session-id")

	if sessionId != "" {
		session, err := c.sessionsClient.GetSessionById("facade", sessionId)
		if session == nil && err == nil {
			err = cerr.NewUnauthorizedError(
				"facade",
				"SESSION_NOT_FOUND",
				"Session invalid or already expired.",
			).WithDetails("session_id", sessionId).WithStatus(440)
		}

		if err == nil {
			// Associate session user with the request
			req = req.WithContext(context.WithValue(req.Context(), "user_id", session.UserId))
			req = req.WithContext(context.WithValue(req.Context(), "user_name", session.UserName))
			req = req.WithContext(context.WithValue(req.Context(), "user", session.User))
			req = req.WithContext(context.WithValue(req.Context(), "session_id", session.Id))
			next.ServeHTTP(res, req)
		} else {
			c.SendError(res, req, err)
		}

	} else {
		next.ServeHTTP(res, req)
	}
}

func (c *SessionsOperationsV1) OpenSession(res http.ResponseWriter, req *http.Request, account *accclients1.AccountV1, roles []string) {
	// let session: SessionV1;
	// let passwordInfo: UserPasswordInfoV1;
	// let settings: ConfigParams;
	// console.log("open session");
	// async.series([
	//     (callback) => {

	//         c.passwordsClient.getPasswordInfo(
	//             nil, account.id, (err, data) => {
	//                 passwordInfo = data;
	//                 callback(err);
	//             }
	//         )
	//     },
	//     // Open a new user session
	//     (callback) => {

	//         let user = <SessionUserV1>{
	//             id: account.id,
	//             name: account.name,
	//             login: account.login,
	//             create_time: account.create_time,
	//             time_zone: account.time_zone,
	//             language: account.language,
	//             theme: account.theme,
	//             roles: roles,
	//             settings: settings,
	//             change_pwd_time: passwordInfo != nil ? passwordInfo.change_time : nil,
	//             custom_hdr: account.custom_hdr,
	//             custom_dat: account.custom_dat
	//         };

	//         let address = HttpRequestDetector.detectAddress(req);
	//         let client = HttpRequestDetector.detectBrowser(req);

	//         c.sessionsClient.openSession(
	//             nil, account.id, account.name,
	//             address, client, user, nil,
	//             (err, data) => {
	//                 session = data;
	//                 callback(err);
	//             }
	//         );
	//     },
	// ], (err) => {
	//     if (err)
	//         c.sendError(req, res, err);
	//     else {
	//         res.json(session);
	//     }
	// });
}

func (c *SessionsOperationsV1) Signup(res http.ResponseWriter, req *http.Request) {
	// let signupData = req.body;
	// let account: AccountV1 = nil;
	// let roles: string[] = signupData.roles != nil && _.isArray(signupData.roles) ? signupData.roles : [];

	// async.series([
	//     // Create account
	//     (callback) => {
	//         let newAccount = <AccountV1>{
	//             name: signupData.name,
	//             login: signupData.login || signupData.email, // Use email by default
	//             language: signupData.language,
	//             theme: signupData.theme,
	//             time_zone: signupData.time_zone
	//         };

	//         c.accountsClient.createAccount(
	//             nil, newAccount,
	//             (err, data) => {
	//                 account = data;
	//                 callback(err);
	//             }
	//         )
	//     },
	//     // Create password for the account
	//     (callback) => {
	//         let password = signupData.password;

	//         c.passwordsClient.setPassword(
	//             nil, account.id, password, callback
	//         );
	//     },
	//     // Create roles for the account
	//     (callback) => {
	//         if (roles.length > 0) {
	//             c.rolesClient.grantRoles(
	//                 nil, account.id, roles, callback
	//             );
	//         } else {
	//             callback();
	//         }
	//     }
	// ], (err) => {
	//     if (err)
	//         c.sendError(req, res, err);
	//     else
	//         c.openSession(req, res, account, roles);
	// });
}

func (c *SessionsOperationsV1) Signin(res http.ResponseWriter, req *http.Request) {
	// let login = req.param("login");
	// let password = req.param("password");

	// let account: AccountV1;
	// let roles: string[] = [];

	// async.series([
	//     // Find user account
	//     (callback) => {
	//         c.accountsClient.getAccountByIdOrLogin(nil, login, (err, data) => {
	//             if (err == nil && data == nil) {
	//                 err = new BadRequestException(
	//                     nil,
	//                     "WRONG_LOGIN",
	//                     "Account " + login + " was not found"
	//                 ).withDetails("login", login);
	//             }

	//             account = data;
	//             callback(err);
	//         });
	//     },
	//     // Authenticate user
	//     (callback) => {
	//         c.passwordsClient.authenticate(nil, account.id, password, (err, result) => {
	//             // wrong password error is UNKNOWN when use http client
	//             if ( (err == nil && result == false) || (err && err.cause == "Invalid password") )  {
	//                 err = new BadRequestException(
	//                     nil,
	//                     "WRONG_PASSWORD",
	//                     "Wrong password for account " + login
	//                 ).withDetails("login", login);
	//             }

	//             callback(err);
	//         });
	//     },
	//     // Retrieve user roles
	//     (callback) => {
	//         if (c.rolesClient) {
	//             c.rolesClient.getRolesById(nil, account.id, (err, data) => {
	//                 roles = data;
	//                 callback(err);
	//             });
	//         } else {
	//             roles = [];
	//             callback();
	//         }
	//     }
	// ], (err) => {
	//     if (err)
	//         c.sendError(req, res, err);
	//     else
	//         c.openSession(req, res, account, roles);
	// });
}

func (c *SessionsOperationsV1) Signout(res http.ResponseWriter, req *http.Request) {
	// if (req.session_id) {
	//     c.sessionsClient.closeSession(nil, req.session_id, (err, session) => {
	//         if (err) c.sendError(req, res, err);
	//         else res.json(204);
	//     });
	// } else {
	//     res.json(204);
	// }
}
