import nsq
import ast
import toml
import logging

def message_handler(message: nsq.Message):
    message.enable_async()

    config = toml.load('config.toml')
    byte_str_body = message.body
    dict_str_body = byte_str_body.decode('UTF-8')
    messageData = ast.literal_eval(dict_str_body)

    logging.warning(messageData['Content'])
    #-----TODO: Implement Server Callback for receiving Assessment Result-----#

    message.finish()

config = toml.load('config.toml')
nsqdAddress = config['nsq']['host'] + ':' + str(config['nsq']['port'])
nsqTopic = config['nsq']['topic']
nsqChannel = config['nsq']['channel']

r = reader = nsq.Reader(
            topic=nsqTopic, channel=nsqChannel, message_handler=message_handler,
            lookupd_connect_timeout=10, requeue_delay=10, 
            nsqd_tcp_addresses=[nsqdAddress], max_in_flight=5, snappy=False
    )
nsq.run()