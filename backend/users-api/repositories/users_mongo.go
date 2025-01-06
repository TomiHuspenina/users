package users

import (
	"context"
	"fmt"
	"log"

	usersDAO "users-api/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type Mongo struct {
	client     *mongo.Client //cliente que contiene la conexion
	database   string        //nombre de la base de datos
	collection string        //coleccion
}

const (
	connectionURI = "mongodb://%s:%s" // %s es marcador de puesto para el host y el puerto
)

func NewMongo(config MongoConfig) Mongo {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	ctx := context.Background()                                 // para manejar cancelaciones o límites de tiempo en las operaciones.
	uri := fmt.Sprintf(connectionURI, config.Host, config.Port) //Construye la URI de conexión utilizando el host y el puerto.
	cfg := options.Client().ApplyURI(uri).SetAuth(credentials)  //Configura las opciones del cliente de MongoDB, incluyendo la URI y las credenciales de autenticación.

	client, err := mongo.Connect(ctx, cfg)
	if err != nil {
		log.Panicf("error connecting to mongo DB: %v", err)
	}

	return Mongo{
		client:     client,
		database:   config.Database,
		collection: config.Collection,
	}
}

func (repository Mongo) GetUserById(id string) (usersDAO.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return usersDAO.User{}, fmt.Errorf("error converting id to mongo ID: %w", err)
	}

	collection := repository.client.Database(repository.database).Collection(repository.collection)
	result := collection.FindOne(bson.M{"_id": objectID})
	if result.Err() != nil {
		return usersDAO.User{}, fmt.Errorf("error finding document: %w", result.Err())
	}

	var userDAO usersDAO.User
	if err := result.Decode(&userDAO); err != nil {
		return usersDAO.User{}, fmt.Errorf("error decoding result: %w", err)
	}
	return userDAO, nil
}

func (repository Mongo) Login(ctx context.Context, user usersDAO.User) (usersDAO.User, error) {
	collection := repository.client.Database(repository.database).Collection(repository.collection)
	result := collection.FindOne(ctx, bson.M{"user": user.User})
	if result.Err() != nil {
		return usersDAO.User{}, fmt.Errorf("error in login: %w", result.Err())
	}

	var foundUser usersDAO.User
	if err := result.Decode(&foundUser); err != nil {
		return usersDAO.User{}, fmt.Errorf("error decoding result: %w", err)
	}
	return foundUser, nil
}

func (repository Mongo) InsertUser(ctx context.Context, user usersDAO.User) error {
	collection := repository.client.Database(repository.database).Collection(repository.collection)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (repository Mongo) GetUserByName(ctx context.Context, user usersDAO.User) (usersDAO.User, error) {
	collection := repository.client.Database(repository.database).Collection(repository.collection)
	result := collection.FindOne(ctx, bson.M{"user": user.User})
	if result.Err() != nil {
		return usersDAO.User{}, fmt.Errorf("error finding user by name: %w", result.Err())
	}

	var foundUser usersDAO.User
	if err := result.Decode(&foundUser); err != nil {
		return usersDAO.User{}, fmt.Errorf("error decoding result: %w", err)
	}
	return foundUser, nil
}

/*
func (repository Mongo) GetHotelByID(ctx context.Context, id string) (hotelsDAO.Hotel, error) {
	// Get from MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return hotelsDAO.Hotel{}, fmt.Errorf("error converting id to mongo ID: %w", err)
	}
	result := repository.client.Database(repository.database).Collection(repository.collection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {

		return hotelsDAO.Hotel{}, fmt.Errorf("error finding document: %w", result.Err())
	}

	// Convert document to DAO
	var hotelDAO hotelsDAO.Hotel
	if err := result.Decode(&hotelDAO); err != nil {
		return hotelsDAO.Hotel{}, fmt.Errorf("error decoding result: %w", err)
	}
	return hotelDAO, nil
}

func (repository Mongo) GetAllHotels(ctx context.Context) ([]hotelsDAO.Hotel, error) {
	// Get from MongoDB

	result, err := repository.client.Database(repository.database).Collection(repository.collection).Find(ctx, bson.M{})
	if err != nil {
		return []hotelsDAO.Hotel{}, fmt.Errorf("error to get all hotels: %w", err)
	}

	// Convert document to DAO
	var hotels []hotelsDAO.Hotel
	for result.Next(ctx) {
		var hotel hotelsDAO.Hotel
		if err := result.Decode(&hotel); err != nil {
			return nil, fmt.Errorf("error decoding document: %w", err)
		}
		hotels = append(hotels, hotel)
		println("Se capo un hotel")
	}

	return hotels, nil
}

func (repository Mongo) InsertHotel(ctx context.Context, hotel hotelsDAO.Hotel) (string, error) {
	result, err := repository.client.Database(repository.database).Collection(repository.collection).InsertOne(ctx, hotel)
	if err != nil {
		return " ", fmt.Errorf("Error inserting new hotel: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("Error converting id to string")
	}
	fmt.Printf("Inserted hotel with ID: %s\n", objectID.Hex())
	return objectID.Hex(), nil
}

func (repository Mongo) UpdateHotel(ctx context.Context, id string, hotel hotelsDAO.Hotel) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting id to mongo ID: %w", err)
	}

	update := bson.M{}

	if hotel.Name != "" {
		update["name"] = hotel.Name
	}
	if hotel.Address != "" {
		update["address"] = hotel.Address
	}
	if hotel.City != "" {
		update["city"] = hotel.City
	}
	if hotel.State != "" {
		update["state"] = hotel.State
	}
	if hotel.Rating != 0 {
		update["rating"] = hotel.Rating
	}
	if len(hotel.Amenities) > 0 {
		update["amenities"] = hotel.Amenities
	}
	if hotel.Price != 0 {
		update["price"] = hotel.Price
	}
	if hotel.Available_rooms != 0 {
		update["available_rooms"] = hotel.Available_rooms
	}

	if len(update) == 0 {
		return fmt.Errorf("no fields to update for hotel ID %s", hotel.Id)
	}

	filter := bson.M{"_id": objectID}
	result, err := repository.client.Database(repository.database).Collection(repository.collection).UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID %s", hotel.Id)
	}

	return nil
}
*/
