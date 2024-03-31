package micro

type UserService interface {
	GetByID(id int)
}
type UserServiceGen struct {
}

//func (u *UserServiceGen) GetByID(id int) {
//	req := &Request{
//		ServiceName: "UserService",
//		MethodName:  "GetByID",
//		Args:        []any{id},
//	}
//}

type Request struct {
	ServiceName string
	MethodName  string
	Args        []any
}
