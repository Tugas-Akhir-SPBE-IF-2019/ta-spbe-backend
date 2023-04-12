# TF
# this code is for TEXT FINDING - INDIKATOR 1
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
import text_finding.preprocess_dokbaru as dokbaru
import text_finding.preprocess_doklama as doklama
import text_finding.indikator01 as indikator01
import text_finding.indikator02 as indikator02
import text_finding.indikator03 as indikator03
import text_finding.indikator04 as indikator04
import text_finding.indikator05 as indikator05
import text_finding.indikator06 as indikator06
import text_finding.indikator07 as indikator07
import text_finding.indikator08 as indikator08
import text_finding.indikator10 as indikator10
import text_finding.highlight_pdf as highlight_pdf


import nsq, ast, toml, logging, requests, threading


def send_result(message_data, config):
    filename = "./src/static/" + message_data['Filename']

    (instansibaru, judulbaru) = dokbaru.pdfparser(filename)
    (instansilama, judullama) = doklama.pdfparser(filename)

    indicator_number = int(message_data['IndicatorNumber'])
    
    text_proof = ""
    if indicator_number == 1:
        text_proof = indikator01.ceklvl(filename)
    elif indicator_number == 2:
        text_proof = indikator02.ceklvl(filename)
    elif indicator_number == 3:
        text_proof = indikator03.ceklvl(filename)
    elif indicator_number == 4:
        text_proof = indikator04.ceklvl(filename)
    elif indicator_number == 5:
        text_proof = indikator05.ceklvl(filename)
    elif indicator_number == 6:
        text_proof = indikator06.ceklvl(filename)
    elif indicator_number == 7:
        text_proof = indikator07.ceklvl(filename)
    elif indicator_number == 8:
        text_proof = indikator08.ceklvl(filename)
    elif indicator_number == 9:
        text_proof = "Not implemented yet"
    elif indicator_number == 10:
        text_proof = indikator10.ceklvl(filename)
    
    proof_pic_files, proof_pages = highlight_pdf.highlight(filename, text_proof, 'Highlight')

    callback_endpoint = 'http://' + config['server']['host'] + '/assessments/result/callback'
    payload = {
        "user_id": message_data['UserId'],
        "assessment_id": message_data['AssessmentId'],
        "recipient_number": message_data['RecipientNumber'],
        "indicator_assessment_id": message_data['IndicatorAssessmentId'],
        "level": 5,
        "explanation": "berdasarkan data dukung yang diberikan, level yang sesuai adalah level 5",
        "support_data_document_id": message_data['Content'],
        "proof": {
            "text": text_proof,
            "picture_url_list": proof_pic_files,
            "page_list": proof_pages
        }
    }
    requests.post(url = callback_endpoint, json = payload)  

def message_handler(message: nsq.Message):
    message.enable_async()

    config = toml.load('config.toml')
    byte_str_body = message.body
    dict_str_body = byte_str_body.decode('UTF-8')
    message_data = ast.literal_eval(dict_str_body)

    logging.warning(message_data['Content'])

    #send result
    send_result(message_data, config)

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