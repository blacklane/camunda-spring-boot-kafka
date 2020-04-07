package com.camunda.demo.springboot.events;

public class RideEvent {
    public RideEvent(String event, String payload) {
        this.event = event;
        this.payload = payload;
    }

    public RideEvent() {
    }

    public String getEvent() {
        return event;
    }

    public void setEvent(String event) {
        this.event = event;
    }

    private String event;

    public String getPayload() {
        return payload;
    }

    public void setPayload(String payload) {
        this.payload = payload;
    }

    private String payload;
}
