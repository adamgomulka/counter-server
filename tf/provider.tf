provider "google" {
    credentials = "${file("./nyt-adam-gomulka-interview-6b516e0cfe83.json")}"
    project = "nyt-adam-gomulka-interview"
    region = "us-east4"
}
