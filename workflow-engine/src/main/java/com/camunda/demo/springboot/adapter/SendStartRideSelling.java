package com.camunda.demo.springboot.adapter;

import com.camunda.demo.springboot.events.RideEvent;
import org.camunda.bpm.engine.delegate.DelegateExecution;
import org.camunda.bpm.engine.delegate.JavaDelegate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Component;

import com.camunda.demo.springboot.ProcessConstants;
import org.springframework.util.concurrent.ListenableFuture;
import org.springframework.util.concurrent.ListenableFutureCallback;

@Component
public class SendStartRideSelling implements JavaDelegate {

  @Autowired
  private KafkaTemplate<String, RideEvent> kafkaTemplate;

  @Value(value = "${kafka.topicSelling}")
  private String topic;

  public void sendKafkaMessage(String msg) {
    RideEvent event = new RideEvent("StartRideSelling", msg);
    ListenableFuture<SendResult<String, RideEvent>> future = kafkaTemplate.send(topic, event);

    future.addCallback(new ListenableFutureCallback<SendResult<String, RideEvent>>() {

      @Override
      public void onSuccess(SendResult<String, RideEvent> result) {
        System.out.println("Sent message=[" + msg +
                "] with offset=[" + result.getRecordMetadata().offset() + "]");
      }
      @Override
      public void onFailure(Throwable ex) {
        System.out.println("Unable to send message=["
                + msg + "] due to : " + ex.getMessage());
      }
    });
  }
  
  @Override
  public void execute(DelegateExecution ctx) throws Exception {
    String orderId = (String) ctx.getVariable(ProcessConstants.VAR_NAME_orderId);
    
    sendKafkaMessage(orderId);
  }

}
