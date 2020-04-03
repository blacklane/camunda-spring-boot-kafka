package com.camunda.demo.springboot.adapter;

import com.camunda.demo.springboot.ProcessConstants;
import org.camunda.bpm.engine.ProcessEngine;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.util.UUID;

@Component
public class KafkaConsumer {

    @Autowired
    private ProcessEngine camunda;

    @KafkaListener(topics = "rides", groupId = "camunda-consumer") //"${kafka.consumerGroupId}")
    public void listen(String message) {
        System.out.println("Received Messasge in group foo: " + message);
        // and call back directly with a generated transactionId
        handleGoodsShippedEvent(message, UUID.randomUUID().toString());
    }



    private void handleGoodsShippedEvent(String orderId, String shipmentId) {
        camunda.getRuntimeService().createMessageCorrelation(ProcessConstants.MSG_NAME_GoodsShipped) //
                .processInstanceVariableEquals(ProcessConstants.VAR_NAME_orderId, orderId) //
                .setVariable(ProcessConstants.VAR_NAME_shipmentId, shipmentId) //
                .correlateWithResult();
    }
}
