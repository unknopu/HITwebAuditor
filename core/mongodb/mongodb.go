package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db *mongo.Database

	// ErrorNotFound error not found
	ErrorNotFound = errors.New("Not found")

	// ErrorInvalidID error invalid id
	ErrorInvalidID = errors.New("Invalid ID")

	// ErrorDucumentDuplicate error ducument is duplicate
	ErrorDucumentDuplicate = errors.New("Ducument is duplicate")
)

// Options mongo option
type Options struct {
	URL          string
	Port         int
	DatabaseName string
	Username     string
	Password     string
	Debug        bool
	Root         bool
}

// InitDatabase new database
func InitDatabase(o *Options) error {
	ctx, concel := context.WithTimeout(context.Background(), 10*time.Second)
	defer concel()

	uri := fmt.Sprintf("mongodb://%s:%d", o.URL, o.Port)
	if o.Username != "" && o.Password != "" {
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
			o.Username, o.Password, o.URL, o.DatabaseName,
		)
	}
	if o.DatabaseName == "" {
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s:%d",
			o.Username, o.Password, o.URL, o.Port,
		)
	}
	log.Println(uri)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	db = client.Database(o.DatabaseName)
	return nil
}

// DB database
func DB() *mongo.Database {
	return db
}

// Repo common repo
type Repo struct {
	Collection *mongo.Collection
	Mux        sync.Mutex
}

// Create create
func (r *Repo) Create(i interface{}) error {
	rt := reflect.TypeOf(i)
	print(rt.Kind())
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return r.createMany(i)
	default:
		return r.createOne(i)
	}
}

func (r *Repo) createOne(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if m, ok := i.(ModelInterface); ok {
		m.Stamp()
		m.SetID(primitive.NewObjectID())
	}
	if _, err := r.Collection.InsertOne(ctx, i); err != nil {
		return WrapError(err)
	}
	return nil
}

// UpdateOneByPrimitiveM update one by primitive m
func (r *Repo) UpdateOneByPrimitiveM(id primitive.ObjectID, u primitive.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r.Mux.Lock()
	_, err := r.Collection.UpdateOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		}, u)
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) createMany(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	v := reflect.ValueOf(i)
	is := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		if m, ok := v.Index(i).Interface().(ModelInterface); ok {
			m.Stamp()
			m.SetID(primitive.NewObjectID())
			is = append(is, m)
		}
	}
	if _, err := r.Collection.InsertMany(ctx, is); err != nil {
		return WrapError(err)
	}
	return nil
}

// Update update
func (r *Repo) Update(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(ModelInterface); ok {
		m.UpdateStamp()
		id = m.GetID()
	}
	r.Mux.Lock()
	err := r.Collection.FindOneAndReplace(ctx, primitive.M{"_id": id}, i).Err()
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// Delete delete
func (r *Repo) Delete(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(ModelInterface); ok {
		m.UpdateDeletedStamp()
		id = m.GetID()
	}
	r.Mux.Lock()
	_, err := r.Collection.UpdateOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		}, primitive.M{
			"$set": i,
		})
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// HardDelete delete
func (r *Repo) HardDelete(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(ModelInterface); ok {
		id = m.GetID()
	}
	r.Mux.Lock()
	_, err := r.Collection.DeleteOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		})
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// FindOneByPrimitiveM find one by primitive.M
func (r *Repo) FindOneByPrimitiveM(m primitive.M, i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := r.Collection.FindOne(ctx, m).Decode(i)
	if err != nil {
		return ErrorNotFound
	}
	return nil
}

// FindOneByID find one by id
func (r *Repo) FindOneByID(id string, i interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrorInvalidID
	}

	err = r.FindOneByPrimitiveM(primitive.M{
		"_id": oid,
	}, i)
	if err != nil {
		return err
	}
	return nil
}

// FindAllByPrimitiveM find all by primitive.M
func (r *Repo) FindAllByPrimitiveM(m primitive.M, result interface{}, opts ...*options.FindOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.Collection.Find(ctx, m, opts...)
	if err != nil {
		return WrapError(err)
	}
	defer func() { _ = cur.Close(ctx) }()
	if err := r.BindModels(ctx, cur, result); err != nil {
		return err
	}
	return nil
}

// BindModels bind array model
func (r *Repo) BindModels(ctx context.Context, cur *mongo.Cursor, result interface{}) error {
	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	i := 0

	for {
		elemp := reflect.New(elemt)
		if !cur.Next(ctx) {
			break
		}
		err := cur.Decode(elemp.Interface())
		if err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		i++
	}
	resultv.Elem().Set(slicev.Slice(0, i))
	return nil
}

// WrapError weap mongo error
func WrapError(err error) error {
	if e, ok := err.(mongo.WriteException); ok {
		if len(e.WriteErrors) > 0 {
			we := e.WriteErrors[0]
			switch we.Code {
			case 11000:
				return ErrorDucumentDuplicate
			default:
				return fmt.Errorf("(%d) %s", we.Code, we.Message)
			}
		}
	}
	return nil
}

// GetLookup Get Look Up
func (r *Repo) GetLookup(collection string, localField string, foreignID string, as string) primitive.M {
	return primitive.M{
		"$lookup": primitive.M{
			"from":         collection,
			"localField":   localField,
			"foreignField": foreignID,
			"as":           as,
		},
	}
}

// GetUnwind Get unwind
func (r *Repo) GetUnwind(path string, preserve bool) primitive.M {
	return primitive.M{
		"$unwind": primitive.M{
			"path":                       path,
			"preserveNullAndEmptyArrays": preserve,
		},
	}
}

// GetCount Get Count
func (r *Repo) GetCount(as string) primitive.M {
	return primitive.M{
		"$count": as,
	}
}

// AggregateAllByPrimitiveA aggregate with pipeline by using primitive A
func (r *Repo) AggregateAllByPrimitiveA(p primitive.A, result interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // TODO: Just for tester to test other issue
	defer cancel()
	opts := options.Aggregate()
	cur, err := r.Collection.Aggregate(ctx, p, opts)
	if err != nil {
		return err
	}
	defer func() {
		cerr := cur.Close(ctx)
		if err == nil {
			err = cerr
		}
	}()

	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	i := 0

	for {
		elemp := reflect.New(elemt)
		if !cur.Next(ctx) {
			break
		}

		if err = cur.Decode(elemp.Interface()); err != nil {
			return err
		}

		slicev = reflect.Append(slicev, elemp.Elem())
		i++
	}

	if err = cur.Err(); err != nil {
		return err
	}

	resultv.Elem().Set(slicev.Slice(0, i))
	return nil
}

// AggregateAllByPrimitiveACount aggregate with pipeline by using primitive A Count
func (r *Repo) AggregateAllByPrimitiveACount(aggCount primitive.A) (resultCount *int32, err error) {
	count := []primitive.M{}
	opts := options.Aggregate().SetMaxTime(5 * time.Minute)
	cursor, err := r.Collection.Aggregate(
		context.TODO(),
		aggCount,
		opts)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	if err = cursor.All(context.TODO(), &count); err != nil {
		//	log.Fatal(err)
		return nil, err
	}
	c := count[0]["count"].(int32)
	return &c, nil
}

// ConvertStringToPrimitiveObjectIDs convert []string to  []primitive.ObjectID
func (r *Repo) ConvertStringToPrimitiveObjectIDs(ids []string) []primitive.ObjectID {
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objectIds = append(objectIds, oid)
		}
	}
	return objectIds
}
