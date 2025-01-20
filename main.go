package ez

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/bytedance/gopkg/lang/fastrand"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
	"bufio"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	_ = bytes.Buffer{}
	_ = fmt.Sprintf("")
	_ = ioutil.ReadFile
	_ = http.Client{}
	_ = json.Marshal
	_ = fastrand.Uint32()
	_ = bcrypt.GenerateFromPassword
	_ = os.Getenv("")
	_ = sync.Mutex{}
	_ = bufio.NewReader

	_ = bson.M{}
	_ = primitive.ObjectID{}
	_ = mongo.Client{}
	_ = options.ClientOptions{}
	_ = readpref.Primary()
	_ = context.TODO()
)

var MongoClient *mongo.Client
var clientLock sync.Mutex
var once sync.Once

func GetMongoClient(URI string) *mongo.Client {
	clientLock.Lock()
	defer clientLock.Unlock()
	if MongoClient == nil || !IsMongoConnected(MongoClient) {
		fmt.Println("MongoDB client is nil or disconnected. Reconnecting...")
		once.Do(func() {
			var err error
			MongoClient, err = connectToMongo(URI)
			if err != nil {
				fmt.Println("Failed to connect to MongoDB: %v", err)
				return
			}
		})
		if !IsMongoConnected(MongoClient) {
			MongoClient.Disconnect(context.TODO())
			MongoClient, _ = connectToMongo(URI)
		}
	}
	return MongoClient
}

func connectToMongo(URI string) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func IsMongoConnected(client *mongo.Client) bool {
	err := client.Ping(context.TODO(), nil)
	return err == nil
}

func Mongoupdate_one(client *mongo.Client, mydb string, mycollection string, filter bson.D, update bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func Mongoupdate_many(client *mongo.Client, mydb string, mycollection string, filter bson.D, update bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func Mongofind_one(client *mongo.Client, mydb string, mycollection string, filter bson.D) (map[string]interface{}, error) {
	collection := client.Database(mydb).Collection(mycollection)
	var result map[string]interface{}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Mongofind_many(client *mongo.Client, mydb string, mycollection string, filter bson.D) ([]map[string]interface{}, error) {
	collection := client.Database(mydb).Collection(mycollection)
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []map[string]interface{}
	for cur.Next(context.Background()) {
		var result map[string]interface{}
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func Mongodel_one(client *mongo.Client, mydb string, mycollection string, filter bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func Mongodel_many(client *mongo.Client, mydb string, mycollection string, filter bson.D) error {
	collection := client.Database(mydb).Collection(mycollection)
	_, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}


func Hash(input string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("hvcjsfhavsfvsa"+input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func Comparehash(in1, in2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(in1), []byte("hvcjsfhavsfvsa"+in2))
	return err == nil
}

func Randint(min, max int) int {
	return min + int(fastrand.Uint32n(uint32(max-min+1)))
}

func Randint64(min, max int) int {
	return min + int(fastrand.Uint64n(uint64(max-min+1)))
}

func Randfloat(min, max float32) float32 {
	return min + fastrand.Float32()*(max-min)
}


func Randfloat64(min, max float64) float64 {
	return min + fastrand.Float64()*(max-min)
}


func Reverseslice(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func InIarray(arr []int, num int) bool {
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

func InSarray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}



func Readfile(filename string, args ...bool) (interface{}, error) {

	Linebyline := false
	if len(args) > 0 {
		Linebyline = args[0]
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if Linebyline {

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading lines: %w", err)
		}
		return lines, nil
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return string(content), nil
}

func WriteFile(filename string, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}


func AppendFile(filename string, content string, args ...bool) error {

	top := false
	addnewline := false
	if len(args) > 0 {
		top = args[0]
	}
	if len(args) > 1 {
		addnewline = args[1]
	}
	
	var existingContent string
	data, err := Readfile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	existingContent = data.(string)

	var newContent string
	if top {
		newContent = content + existingContent
	} else {
		newContent = existingContent + content
	}
	if addnewline {
		newContent += "\n"
	}
	err = ioutil.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}


func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func executeRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func ParseJson(body string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return result, nil
}

func Get(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating GET request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Post(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}



func Put(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating PUT request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Delete(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating DELETE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Patch(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating PATCH request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Options(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating OPTIONS request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Head(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HEAD request: %w", err)
	}
	addHeaders(req, headers)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing HEAD request: %w", err)
	}
	defer resp.Body.Close()

	return resp.Status, nil
}


func Trace(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("TRACE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating TRACE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}


func GetJson(url string, headers map[string]string) (map[string]interface{}, error) {
	body, err := Get(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing GET request: %w", err)
	}
	return ParseJson(body)
}

func PostJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Post(url, data, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing POST request: %w", err)
	}
	return ParseJson(Body)
}

func PutJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Put(url, data, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing PUT request: %w", err)
	}
	return ParseJson(Body)
}

func DeleteJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Delete(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing DELETE request: %w", err)
	}
	return ParseJson(Body)
}

func PatchJson(url string, data []byte, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Patch(url, data, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing PATCH request: %w", err)
	}
	return ParseJson(Body)
}

func OptionsJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Options(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing OPTIONS request: %w", err)
	}
	return ParseJson(Body)
}

func HeadJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Head(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing HEAD request: %w", err)
	}
	return ParseJson(Body)
}

func TraceJson(url string, headers map[string]string) (map[string]interface{}, error) {
	Body, err := Trace(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error performing TRACE request: %w", err)
	}
	return ParseJson(Body)
}
