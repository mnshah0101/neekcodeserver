package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSolution(w http.ResponseWriter, r *http.Request) {

	question_id := r.FormValue("question_id")

	if question_id == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")

	if date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Upload file to S3
	fileName := header.Filename
	fileURL, err := UploadFileToS3(fileName, fileContent)

	if err != nil {

		log.Print("Error uploading file: ", err.Error())

		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	// Save solution to MongoDB

	collection := MongoClient.Database("test").Collection("solutions")

	solution := Solution{
		Username:    username,
		Email:       email,
		File_url:    fileURL,
		Likes:       0,
		Question_id: question_id,
		Date:        date,
	}

	_, err = collection.InsertOne(context.Background(), solution)

	if err != nil {
		http.Error(w, "Error saving solution", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Solution created successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GetSolutions(w http.ResponseWriter, r *http.Request) {

	question_id := r.URL.Query().Get("question_id")

	if question_id == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	collection := MongoClient.Database("test").Collection("solutions")

	cursor, err := collection.Find(context.Background(), bson.M{"question_id": question_id})

	if err != nil {
		http.Error(w, "Error fetching solutions", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(context.Background())

	var solutions []Solution

	for cursor.Next(context.Background()) {
		var solution Solution
		cursor.Decode(&solution)
		solutions = append(solutions, solution)
	}

	response := map[string]interface{}{
		"solutions": solutions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func LikeSolution(w http.ResponseWriter, r *http.Request) {

	solution_id := r.FormValue("solution_id")

	if solution_id == "" {
		http.Error(w, "Solution ID is required", http.StatusBadRequest)
		return
	}

	collection := MongoClient.Database("test").Collection("solutions")

	solutionID, err := primitive.ObjectIDFromHex(solution_id)
	if err != nil {
		log.Println("Invalid solution ID format:", err)
		http.Error(w, "Invalid solution ID format", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": solutionID}

	update := bson.M{"$inc": bson.M{"likes": 1}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		http.Error(w, "Error liking solution", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Solution liked successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func DislikeSolution(w http.ResponseWriter, r *http.Request) {
	solution_id := r.FormValue("solution_id")

	if solution_id == "" {
		http.Error(w, "Solution ID is required", http.StatusBadRequest)
		return
	}

	collection := MongoClient.Database("test").Collection("solutions")

	solutionID, err := primitive.ObjectIDFromHex(solution_id)
	if err != nil {
		log.Println("Invalid solution ID format:", err)
		http.Error(w, "Invalid solution ID format", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": solutionID}

	update := bson.M{"$inc": bson.M{"likes": -1}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		http.Error(w, "Error liking solution", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Solution liked successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
