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
import text_similarity.preprocess as ts_preprocess
import random
import pickle
import pandas as pd

import nsq, ast, toml, logging, requests, threading

# TODO model 2 and 10 are still missing
model_dt_1 = pickle.load(open('./decision_tree_model/model_dtc_1.pckl', 'rb'))
model_dt_3 = pickle.load(open('./decision_tree_model/model_dtc_3.pckl', 'rb'))
model_dt_4 = pickle.load(open('./decision_tree_model/model_dtc_4.pckl', 'rb'))
model_dt_5 = pickle.load(open('./decision_tree_model/model_dtc_5.pckl', 'rb'))
model_dt_6 = pickle.load(open('./decision_tree_model/model_dtc_6.pckl', 'rb'))
model_dt_7 = pickle.load(open('./decision_tree_model/model_dtc_7.pckl', 'rb'))
model_dt_8 = pickle.load(open('./decision_tree_model/model_dtc_8.pckl', 'rb'))
model_dt_9 = pickle.load(open('./decision_tree_model/model_dtc_9.pckl', 'rb'))


def send_result(message_data, config):
    new_document_list = []
    old_document_list = []
    meeting_minutes_list = []
    for support_document in message_data['support_document_list']:
        if support_document['type'] == "NEW_DOCUMENT":
            new_document_list.append(support_document)
        elif support_document['type'] == "OLD_DOCUMENT":
            old_document_list.append(support_document)
        else: # support_document['type'] == "MEETING_MINUTES"
            meeting_minutes_list.append(support_document)

    new_document_title_list = []
    for new_document in new_document_list:
        filename = "./src/static/" + new_document['name']
        original_filename = new_document['original_name']

        (instansibaru, judulbaru) = dokbaru.pdfparser(filename)
        new_document_title_list.append(judulbaru)

    # indicator_number = int(message_data['IndicatorNumber'])
    # indicator_detail = (message_data['IndicatorDetail'])
    # institution_name = (message_data['InstitutionName'])

    proof_list = []
    for indicator_asssessment in message_data['indicator_assessment_list']:
        indicator_number = indicator_asssessment['number']
        document_proof_list = []

        for idx, new_document in enumerate(new_document_list):
            filename = "./src/static/" + new_document['name']
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
                # TODO not yet implemented
                text_proof = ""
            elif indicator_number == 10:
                text_proof = indikator10.ceklvl(filename)
            
            if text_proof != "":
                proof_pic_files, proof_pages, page_with_matches = highlight_pdf.highlight(filename, text_proof, 'Highlight')
                document_proof_list.append(
                    {
                        "name": new_document['name'],
                        "original_name": new_document['original_name'],
                        "type": "NEW_DOCUMENT",
                        "text": text_proof,
                        "title": new_document_title_list[idx],
                        "picture_file_list": proof_pic_files,
                        "specific_page_document_url": proof_pages,
                        "document_page_list": page_with_matches
                    }
                )

        proof_list.append(
            {
                "indicator_assessment": indicator_asssessment,
                "document_proof": document_proof_list
            }
        )
    logging.warning(proof_list)

    result_list = []
    for proof in proof_list:
        document_proof_list = []
        page_text = ""
        explanation_text = ""
        level = 0 # init value
        for document_proof in proof['document_proof']:
            text_document_proof = document_proof['text']
            # TODO integrate with model; loop all of the old and meeting minutes documents
            # level = random.randint(2,4)
            # document_proof_list.append(
            #     {
            #         'name': document_proof['name'],
            #         'original_name': document_proof['original_name'],
            #         'text': document_proof['text'],
            #         "picture_file_list": document_proof['picture_file_list'],
            #         "specific_page_document_url": document_proof['specific_page_document_url'],
            #         "document_page_list": document_proof['document_page_list']
            #     }
            # )
            for idx, page in enumerate(document_proof['document_page_list']):
                if idx == 0:
                    page_text += str(page)
                elif idx < len(proof_pages) - 1:
                    page_text += f", {page}"
                else:
                    if idx == 1:
                        page_text += f" dan {page}"
                    else:
                        page_text += f", dan {page}"
            original_filename = document_proof['original_name']
            judulbaru = document_proof['title']
            institution_name = message_data['institution_name']
            indicator_detail = proof['indicator_assessment']['detail']
            indicator_number = proof['indicator_assessment']['number']


            # TODO placeholder value
            judul_notulensi_undangan = '[]'
            isi_notulensi_undangan = '[]'
            JudulDokumenLama = '[]'
            TeksDokumenLama = '[]'
            data_prediction = {
                'Judul': [judulbaru],
                'teks': [text_document_proof],
                'JudulDokumenLama': [JudulDokumenLama],
                'TeksDokumenLama': [TeksDokumenLama],
                'JudulNotulensiUndangan': [judul_notulensi_undangan],
                'IsiNotulensiUndangan': [isi_notulensi_undangan]
            }
            df_predict = pd.DataFrame(data_prediction)
            if indicator_number == 1:
                logging.warning("prediksi indikator 1")
                level = ts_preprocess.predict_model(model_dt_1, df_predict)[0]
            elif indicator_number == 2:
                logging.warning("prediksi indikator 2")
            elif indicator_number == 3:
                logging.warning("prediksi indikator 3")
                level = ts_preprocess.predict_model(model_dt_3, df_predict)[0]
            elif indicator_number == 4:
                logging.warning("prediksi indikator 4")
                level = ts_preprocess.predict_model(model_dt_4, df_predict)[0]
            elif indicator_number == 5:
                logging.warning("prediksi indikator 5")
                level = ts_preprocess.predict_model(model_dt_5, df_predict)[0]
            elif indicator_number == 6:
                logging.warning("prediksi indikator 6")
                level = ts_preprocess.predict_model(model_dt_6, df_predict)[0]
            elif indicator_number == 7:
                logging.warning("prediksi indikator 7")
                level = ts_preprocess.predict_model(model_dt_7, df_predict)[0]
            elif indicator_number == 8:
                logging.warning("prediksi indikator 8")
                level = ts_preprocess.predict_model(model_dt_8, df_predict)[0]
            elif indicator_number == 9:
                logging.warning("prediksi indikator 9")
                level = ts_preprocess.predict_model(model_dt_9, df_predict)[0]
            elif indicator_number == 10:
                logging.warning("prediksi indikator 10")

            explanation_text = f"Verifikasi dan validasi telah dilakukan terhadap penjelasan dan data dukung pada Indikator {indicator_number} {indicator_detail} pada {institution_name}, dimana tercantum dalam {judulbaru}, yaitu pada halaman {page_text} sesuai data dukung  {original_filename}"
        result_list.append({
            'indicator_assessment': proof['indicator_assessment'],
            'document_proof': proof['document_proof'],
            'result': {
                'level': int(level),
                'explanation': explanation_text
            }
        })
    logging.warning(result_list)
    callback_endpoint = 'http://' + config['server']['host'] + '/assessments/result/callback'
    payload = {
        "user_id": message_data['user_id'],
        "assessment_id": message_data['assessment_id'],
        "recipient_number": message_data['recipient_number'],
        "data": result_list
    }
    logging.warning(payload)
    requests.post(url = callback_endpoint, json = payload)  


def message_handler(message: nsq.Message):
    message.enable_async()

    config = toml.load('config.toml')
    byte_str_body = message.body
    dict_str_body = byte_str_body.decode('UTF-8')
    message_data = ast.literal_eval(dict_str_body)

    # logging.warning(message_data['Content'])
    logging.warning(message_data['indicator_assessment_list'])

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