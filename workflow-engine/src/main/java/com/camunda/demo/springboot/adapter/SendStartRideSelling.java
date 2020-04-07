package com.camunda.demo.springboot.adapter;

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
  private KafkaTemplate<String, String> kafkaTemplate;

  @Value(value = "${kafka.topicRides}")
  private String topic;

  public void sendKafkaMessage(String msg) {
    ListenableFuture<SendResult<String, String>> future = kafkaTemplate.send(topic, msg);

    future.addCallback(new ListenableFutureCallback<SendResult<String, String>>() {

      @Override
      public void onSuccess(SendResult<String, String> result) {
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
