#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>

const char *WIFI_SSID = "Your SSID";
const char *WIFI_PASSWORD = "Your Password";

const char *URL = "https://hooks.slack.com/services/SECRET";
const char *SPACEAPITOKEN = "SECRET";


WiFiClientSecure client;
HTTPClient httpsClient;

int lastState = 0;
int currentState = 0;

void setup() {
  Serial.begin(9600);

  WiFi.mode(WIFI_STA);

  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("Connected");

  client.setInsecure();

  pinMode(D8, INPUT);
  pinMode(D6, OUTPUT);
}

void loop() {
  if (digitalRead(D8) == LOW) {
    currentState = 1;
  } else {
    currentState = 0;
  }

  if (lastState != currentState) {
    if (currentState == 1) {
      Serial.print("Send Space ist auf");
      send("Space ist auf!");
      setApiOpen();
    } else {
      Serial.print("Send Space ist geschlossen");
      send("Space ist zu!");
      setApiClose();
    }
  }

  if (digitalRead(D8) == LOW) {
    lastState = 1;
    digitalWrite(D6, HIGH);
  } else {
    lastState = 0;
    digitalWrite(D6, LOW);
  }

  delay(1000);
}

void send(const String &text) {
  String data = "{\"text\":\"" + String(text) + "\"}";
  httpsClient.begin(client, URL);
  httpsClient.addHeader("Content-Type", "application/json");
  httpsClient.POST(data);
  httpsClient.end();
  Serial.print("Send");
}

void setApiClose() {
  httpsClient.begin(client, "https://api.chaostreff-flensburg.de/close");
  httpsClient.POST(SPACEAPITOKEN);
  httpsClient.end();
  Serial.print("Close API");
}

void setApiOpen() {
  httpsClient.begin(client, "https://api.chaostreff-flensburg.de/open");
  httpsClient.POST(SPACEAPITOKEN);
  httpsClient.end();
  Serial.print("OPEN API");
}
