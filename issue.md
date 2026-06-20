# Issue: Registration Endpoint Bugs

## Description
When submitting a registration request with only the `email` and `password` fields, the request succeeds but results in two significant issues:
1. The `name` and `username` fields in the database are saved as empty strings, even though they are marked as required in the `request.RegisterRequest` struct.
2. The JSON response returns the newly created user object, which exposes the hashed password.

## Problem 1: Missing Request Validation
The `request.RegisterRequest` struct in `internal/delivery/http/request/auth.go` properly defines `validate:"required"` tags for `Username` and `Name`. However, the validation is never actually executed.

In `internal/service/auth.go`, the `AuthService` struct has a `validate *validator.Validate` instance injected into it, but `s.validate.Struct(req)` is never called inside the `Register` method. Because of this, the empty string values from the JSON binding are passed directly to the database.

### Proposed Fix for Problem 1
Modify `internal/service/auth.go` within the `Register` method to validate the request before proceeding:
```go
func (s *AuthService) Register(req request.RegisterRequest) (entity.User, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.User{}, err
	}

	// ... continue with password hashing and user creation
}
```

## Problem 2: Exposed Password Hash in Response
In `internal/delivery/handler/auth.go`, the `Register` handler directly returns the `result` from `h.IAuthService.Register(req)` inside the `Data` field of the JSON response. Since `result` is of type `entity.User`, and `entity.User` defines `json:"password"`, the hash is leaked.

In other handlers like `Login` and `Refresh`, the code maps the user object to a safe `response.UserInfo` struct to prevent exposing sensitive fields.

### Proposed Fix for Problem 2
Modify `internal/delivery/handler/auth.go` within the `Register` method to map the `result` to a `response.UserInfo` struct before returning it:
```go
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "user registered successfully",
		Data: response.UserInfo{
			Id:              result.Id,
			Username:        result.Username,
			Name:            result.Name,
			Email:           result.Email,
			Status:          result.Status,
			EmailVerifiedAt: result.EmailVerifiedAt.Format(time.RFC3339),
		},
	})
```
Alternatively, update `internal/entity/user.go` to use `json:"-"` for the `Password` field so it is never accidentally serialized into JSON responses. Using `response.UserInfo` is the recommended immediate fix to stay consistent with the `Login` handler.
