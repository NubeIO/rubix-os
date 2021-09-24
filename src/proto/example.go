package proto

// GOLANG
/*
INSTALL
go install google.golang.org/protobuf/cmd/protoc-gen-go


BUILD
protoc -I . --go_out=. ./person.proto


RUN
go run *.go
*/

//PROTO FILE
/*
make the proto file
"""

syntax="proto3";

message Person {
      string name = 1;
      int32 age = 2;
}
"""

//PYTHON
"""


install
pip3 install protobuf

build
protoc -I=./ --python_out=./ person.proto

"""


# run the python file
"""
import person_pb2 as p

aa = p.Person = "aaaa"
print(p.Person)

"""
*/

/*



func main() {

	elliot := &Person{
		Name: "Elliot",
		Age:  24,
	}

	data, err := proto.Marshal(elliot)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	// printing out our raw protobuf object
	fmt.Println(data)

	// let's go the other way and unmarshal
	// our byte array into an object we can modify
	// and use
	newElliot := &Person{}
	err = proto.Unmarshal(data, newElliot)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	// print out our `newElliot` object
	// for good measure
	fmt.Println(newElliot.GetAge())
	fmt.Println(newElliot.GetName())

}


*/
