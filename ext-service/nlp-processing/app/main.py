import nsq, ast, toml, logging, requests, threading

def send_result(message_data, config):
    callback_endpoint = 'http://' + config['server']['host'] + '/assessments/result/callback'
    payload = {
        "user_id": message_data['UserId'],
        "assessment_id": message_data['AssessmentId'],
        "recipient_number": message_data['RecipientNumber'],
        "indicator_assessment_id": message_data['IndicatorAssessmentId'],
        "level": 4,
        "explanation": "berdasarkan data dukung yang diberikan, level yang sesuai adalah level 4",
        "support_data_document_id": message_data['Content'],
        "proof": "<p>ini <b>buktinya</b></p>"
        
    }
    requests.post(url = callback_endpoint, json = payload)  

def message_handler(message: nsq.Message):
    message.enable_async()

    config = toml.load('config.toml')
    byte_str_body = message.body
    dict_str_body = byte_str_body.decode('UTF-8')
    message_data = ast.literal_eval(dict_str_body)

    logging.warning(message_data['Content'])
    #-----TODO: Implement Server Callback for receiving Assessment Result-----#
    timer = threading.Timer(config['callback']['mockcooldown'], send_result, args=(message_data, config))
    timer.start()

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