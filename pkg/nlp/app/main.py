import nsq, ast, toml, logging, requests
import text_similarity.process_similarity as similarity



def message_handler(message: nsq.Message):
    message.enable_async()

    config = toml.load('config.toml')
    byte_str_body = message.body
    dict_str_body = byte_str_body.decode('UTF-8')
    message_data = ast.literal_eval(dict_str_body)

    logging.warning(message_data['indicator_assessment_list'])

    #send result
    similarity_result_payload = similarity.process_similarity(message_data)
    callback_endpoint = 'http://' + config['server']['host'] + '/assessments/result/callback'
    requests.post(url = callback_endpoint, json = similarity_result_payload)  

    message.finish()

config = toml.load('config.toml')
nsqd_address = config['nsq']['host'] + ':' + str(config['nsq']['port'])
nsq_topic = config['nsq']['topic']
nsq_channel = config['nsq']['channel']

r = reader = nsq.Reader(
            topic=nsq_topic, channel=nsq_channel, message_handler=message_handler,
            lookupd_connect_timeout=10, requeue_delay=10, 
            nsqd_tcp_addresses=[nsqd_address], max_in_flight=5, snappy=False
    )
nsq.run()