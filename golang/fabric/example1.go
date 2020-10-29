package fabric

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	shim "github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"io"
	"strconv"
	"time"
)

type SimpleChaincode struct {
}

var BackGroundNo int = 0

var RecordNo int = 0

type School struct {
	Name           string
	Location       string
	Address        string
	PriKey         string
	PubKey         string
	StudentAddress []string
}

type Student struct {
	Name         string
	Address      string
	BackgroundId []int
}

type Background struct {
	Id       int
	ExitTime int64
	Status   string
}

type Record struct {
	Id              int
	SchoolAddress   string
	StudentAddress  string
	SchoolSign      string
	ModifyTime      int64
	ModifyOperation string
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

/*
 * 客户端发起请求执行"diploma"智能合约时会调用 Invoke 方法
 */
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// 获取请求调用智能合约的方法和参数
	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createSchool" {
		return t.createSchool(stub, args)
	} else if function == "createStudent" {
		return t.createStudent(stub, args)
	} else if function == "enrollStudent" {
		return t.enrollStudent(stub, args)
	} else if function == "updateDiploma" {
		return t.updateDiploma(stub, args)
	} else if function == "getRecords" {
		return t.getRecords(stub)
	} else if function == "getRecordById" {
		return t.getRecordById(stub, args)
	} else if function == "getStudentByAddress" {
		return t.getStudentByAddress(stub, args)
	} else if function == "getSchoolByAddress" {
		return t.getSchoolByAddress(stub, args)
	} else if function == "getBackgroundById" {
		return t.getBackgroundById(stub, args)
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) createSchool(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var school School
	var schoolBytes []byte
	var stuAddress []string
	var address, priKey, pubKey string
	address, priKey, pubKey = GetAddress()

	school = School{Name: args[0], Location: args[1], Address: address, PriKey: priKey, PubKey: pubKey, StudentAddress: stuAddress}
	err := writeSchool(stub, school)
	if err != nil {
		shim.Error("Error write school")
	}

	schoolBytes, err = json.Marshal(&school)
	if err != nil {
		return shim.Error("Error retrieving schoolBytes")
	}

	return shim.Success(schoolBytes)
}

func (t *SimpleChaincode) createStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var student Student
	var studentBytes []byte
	var stuAddress string
	var bgId []int
	stuAddress, _, _ = GetAddress()

	student = Student{Name: args[0], Address: stuAddress, BackgroundId: bgId}
	err := writeStudent(stub, student)
	if err != nil {
		return shim.Error("Error write student")
	}

	studentBytes, err = json.Marshal(&student)
	if err != nil {
		return shim.Error("Error retrieving studentBytes")
	}

	return shim.Success(studentBytes)
}

func (t *SimpleChaincode) getStudentByAddress(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stuBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error("Error retrieving data")
	}
	return shim.Success(stuBytes)
}

func (t *SimpleChaincode) enrollStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	schAddress := args[0]
	schoolSign := args[1]
	stuAddress := args[2]

	var school School
	var schBytes []byte
	var err error

	schBytes, err = stub.GetState(schAddress)
	if err != nil {
		return shim.Error("Error retrieving data")
	}
	err = json.Unmarshal(schBytes, &school)
	if err != nil {
		return shim.Error("Error unmarshalling data")
	}

	var record Record
	record = Record{Id: RecordNo, SchoolAddress: schAddress, StudentAddress: stuAddress, SchoolSign: schoolSign, ModifyTime: time.Now().Unix(), ModifyOperation: "2"}

	err = writeRecord(stub, record)
	if err != nil {
		return shim.Error("Error write record")
	}

	school.StudentAddress = append(school.StudentAddress, stuAddress)
	err = writeSchool(stub, school)
	if err != nil {
		return shim.Error("Error write school")
	}

	RecordNo = RecordNo + 1
	recordBytes, err := json.Marshal(&record)

	if err != nil {
		return shim.Error("Error retrieving recordBytes")
	}

	return shim.Success(recordBytes)

}

func (t *SimpleChaincode) updateDiploma(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	schAddress := args[0]
	schoolSign := args[1]
	stuAddress := args[2]
	modOperation := args[3]

	var recordBytes []byte
	var student Student
	var stuBytes []byte
	var err error

	stuBytes, err = stub.GetState(stuAddress)
	if err != nil {
		return shim.Error("Error retrieving data")
	}
	err = json.Unmarshal(stuBytes, &student)
	if err != nil {
		return shim.Error("Error unmarshalling data")
	}

	var record Record
	record = Record{Id: RecordNo, SchoolAddress: schAddress, SchoolSign: schoolSign, ModifyTime: time.Now().Unix(), ModifyOperation: modOperation}

	err = writeRecord(stub, record)
	if err != nil {
		return shim.Error("Error write record")
	}

	var background Background
	background = Background{Id: BackGroundNo, ExitTime: time.Now().Unix(), Status: modOperation}
	err = writeBackground(stub, background)

	if err != nil {
		return shim.Error("Error write background")
	}

	if modOperation == "0" {
		student.BackgroundId = append(student.BackgroundId, BackGroundNo)
		student = Student{Name: student.Name, Address: student.Address, BackgroundId: student.BackgroundId}
		err = writeStudent(stub, student)
		if err != nil {
			return shim.Error("Error write student")
		}
	}

	BackGroundNo = BackGroundNo + 1
	recordBytes, err = json.Marshal(&record)
	if err != nil {
		return shim.Error("Error retrieving schoolBytes")
	}

	return shim.Success(recordBytes)
}

func (t *SimpleChaincode) getSchoolByAddress(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	schBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error("Error retrieving data")
	}

	return shim.Success(schBytes)
}

func (t *SimpleChaincode) getRecords(stub shim.ChaincodeStubInterface) pb.Response {
	var records []Record
	var number string
	var err error
	var record Record
	var recBytes []byte

	for i := 0; i < RecordNo; i++ {
		number = strconv.Itoa(i)
		recBytes, err = stub.GetState("Record" + number)
		if err != nil {
			return shim.Error("Error get detail")
		}
		err := json.Unmarshal(recBytes, &record)
		if err != nil {
			return shim.Error("Error unmarshalling data")
		}
		records = append(records, record)
		if i == 10 {
			break
		}
	}

	recordBytes, err := json.Marshal(&record)
	if err != nil {
		shim.Error("Error get records")
	}

	return shim.Success(recordBytes)
}

func (t *SimpleChaincode) getBackgroundById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	backBytes, err := stub.GetState("BackGround" + args[0])
	if err != nil {
		return shim.Error("Error retrieving data")
	}
	return shim.Success(backBytes)
}

func writeRecord(stub shim.ChaincodeStubInterface, record Record) error {
	var recID string
	recordBytes, err := json.Marshal(&record)
	if err != nil {
		return err
	}

	recID = strconv.Itoa(record.Id)
	err = stub.PutState("Record"+recID, recordBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func (t *SimpleChaincode) getRecordById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recBytes, err := stub.GetState("Record" + args[0])
	if err != nil {
		return shim.Error("Error retrieving data")
	}

	return shim.Success(recBytes)
}

func writeSchool(stub shim.ChaincodeStubInterface, school School) error {
	schBytes, err := json.Marshal(&school)
	if err != nil {
		return err
	}

	err = stub.PutState(school.Address, schBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeBackground(stub shim.ChaincodeStubInterface, background Background) error {
	var backId string
	backBytes, err := json.Marshal(&background)
	if err != nil {
		return err
	}

	backId = strconv.Itoa(background.Id)
	err = stub.PutState("BackGround"+backId, backBytes)
	if err != nil {
		return errors.New("PutState errpor" + err.Error())
	}

	return nil
}

func writeStudent(stub shim.ChaincodeStubInterface, student Student) error {
	stuBytes, err := json.Marshal(&student)
	if err != nil {
		return err
	}

	err = stub.PutState(student.Address, stuBytes)
	if err != nil {
		return errors.New("PutState error " + err.Error())
	}

	return nil
}

func GetAddress() (string, string, string) {
	var address, priKey, pubKey string

	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", "", ""
	}

	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))

	address = hex.EncodeToString(h.Sum(nil))
	priKey = address + "1"
	pubKey = address + "2"

	return address, priKey, pubKey
}
