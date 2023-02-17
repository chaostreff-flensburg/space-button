#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>

const char *WIFI_SSID = "Chaostreff-Flensburg";
const char *WIFI_PASSWORD = "SECRET";

const char *URL = "https://hooks.slack.com/services/SECRET";
const char *FINGERPRINT = "82AEFD933630DA030A2F6353DE2EB0438BF441F6";

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

  client.setFingerprint(FINGERPRINT);

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
    } else {
      Serial.print("Send Space ist geschlossen");
      send("Space ist zu!");
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
