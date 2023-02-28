#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>

const char *WIFI_SSID = "SSID";
const char *WIFI_PASSWORD = "SECRET";

const char *URL = "https://hooks.slack.com/services/SECRET";
const char *SPACEAPITOKEN = "TOKEN";


WiFiClientSecure client;
HTTPClient httpsClient;

int lastState = 0;
int currentState = 0;

void setup() {
  WiFi.mode(WIFI_STA);

  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
  }

  client.setInsecure();

  pinMode(D8, INPUT);
  pinMode(D6, OUTPUT);

  digitalWrite(D6, HIGH);
  delay(100);
  digitalWrite(D6, LOW);
  delay(100);
  digitalWrite(D6, HIGH);
  delay(100);
  digitalWrite(D6, LOW);
  digitalWrite(D6, HIGH);
  delay(100);
  digitalWrite(D6, LOW);
}

void loop() {
  if (digitalRead(D8) == LOW) {
    currentState = 1;
  } else {
    currentState = 0;
  }

  if (lastState != currentState) {
    if (currentState == 1) {
      send("Space ist auf!");
      setApiOpen();
    } else {
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

  delay(200);
}

void send(const String &text) {
  String data = "{\"text\":\"" + String(text) + "\"}";
  httpsClient.begin(client, URL);
  httpsClient.addHeader("Content-Type", "application/json");
  httpsClient.POST(data);
  httpsClient.end();
}

void setApiClose() {
  httpsClient.begin(client, "https://api.chaostreff-flensburg.de/close");
  httpsClient.POST(SPACEAPITOKEN);
  httpsClient.end();
}

void setApiOpen() {
  httpsClient.begin(client, "https://api.chaostreff-flensburg.de/open");
  httpsClient.POST(SPACEAPITOKEN);
  httpsClient.end();
}
