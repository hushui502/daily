package user

//import (
//	"awesomeProject2/jianyu/unittest/mock"
//	"github.com/golang/mock/gomock"
//	"testing"
//)
//
//func TestUser_GetUserInfo(t *testing.T) {
//	ctr := gomock.NewController(t)
//	defer ctr.Finish()
//
//	var id int64 = 1
//	mockMale := mock.NewMockMale(ctr)
//	gomock.InOrder(
//		mockMale.EXPECT().Get(id).Return(nil),
//		)
//
//	user := NewUser(mockMale)
//	err := user.GetUserInfo(id)
//	if err != nil {
//		t.Errorf("user.GetUserInfo err : %v", err)
//	}
//
//}
