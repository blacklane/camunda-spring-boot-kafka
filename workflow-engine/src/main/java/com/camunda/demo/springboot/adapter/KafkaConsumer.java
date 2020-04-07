package com.camunda.demo.springboot.adapter;

import com.camunda.demo.springboot.ProcessConstants;
import com.camunda.demo.springboot.events.RideEvent;
import org.camunda.bpm.engine.ProcessEngine;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Component;

import java.util.UUID;

@Component
public class KafkaConsumer {

    @Autowired
    private ProcessEngine camunda;

    @KafkaListener(
            topics = "${kafka.topicRides}",
            groupId = "${kafka.consumerGroupId}"
    )
    public void listen(@Payload RideEvent message) {
        System.out.println("Received Message on topic rides: " + message);

        if(message.getEvent().equals("RideAccepted")) {
            // and call back directly with a generated transactionId
            handleRideAcceptedEvent(message.getPayload(), UUID.randomUUID().toString());
        }
    }



    private void handleRideAcceptedEvent(String orderId, String shipmentId) {
        System.out.println("Setting ride to accepted for orderId: " + orderId);
        camunda.getRuntimeService().createMessageCorrelation(ProcessConstants.MSG_NAME_GoodsShipped) //
                .processInstanceVariableEquals(ProcessConstants.VAR_NAME_orderId, orderId) //
                .setVariable(ProcessConstants.VAR_NAME_shipmentId, shipmentId) //
                .correlateWithResult();
    }
}
