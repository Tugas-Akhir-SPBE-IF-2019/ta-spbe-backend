import text_finding.preprocess_dokbaru as text_finding_dokbaru
import text_finding.preprocess_notulen as text_finding_notulen
import text_finding.indikator01 as indikator01
import text_finding.indikator02 as indikator02
import text_finding.indikator03 as indikator03
import text_finding.indikator04 as indikator04
import text_finding.indikator05 as indikator05
import text_finding.indikator06 as indikator06
import text_finding.indikator07 as indikator07
import text_finding.indikator08 as indikator08
import text_finding.indikator09 as indikator09
import text_finding.indikator10 as indikator10
import text_finding.highlight_pdf as highlight_pdf
import text_similarity.preprocess as ts_preprocess
import pickle
import pandas as pd
import logging

model_lsa_svm_1 = pickle.load(open('./lsa_svm_model/lsa_model_svm_1.pkl', 'rb'))
model_lsa_svm_2 = pickle.load(open('./lsa_svm_model/lsa_model_svm_2.pkl', 'rb'))
model_lsa_svm_3 = pickle.load(open('./lsa_svm_model/lsa_model_svm_3.pkl', 'rb'))
model_lsa_svm_4 = pickle.load(open('./lsa_svm_model/lsa_model_svm_4.pkl', 'rb'))
model_lsa_svm_5 = pickle.load(open('./lsa_svm_model/lsa_model_svm_5.pkl', 'rb'))
model_lsa_svm_6 = pickle.load(open('./lsa_svm_model/lsa_model_svm_6.pkl', 'rb'))
model_lsa_svm_7 = pickle.load(open('./lsa_svm_model/lsa_model_svm_7.pkl', 'rb'))
model_lsa_svm_8 = pickle.load(open('./lsa_svm_model/lsa_model_svm_8.pkl', 'rb'))
model_lsa_svm_9 = pickle.load(open('./lsa_svm_model/lsa_model_svm_9.pkl', 'rb'))
model_lsa_svm_10 = pickle.load(open('./lsa_svm_model/lsa_model_svm_10.pkl', 'rb'))

model_svm_1 = pickle.load(open('./lsa_svm_model/model_svm_1.pkl', 'rb'))
model_svm_2 = pickle.load(open('./lsa_svm_model/model_svm_2.pkl', 'rb'))
model_svm_3 = pickle.load(open('./lsa_svm_model/model_svm_3.pkl', 'rb'))
model_svm_4 = pickle.load(open('./lsa_svm_model/model_svm_4.pkl', 'rb'))
model_svm_5 = pickle.load(open('./lsa_svm_model/model_svm_5.pkl', 'rb'))
model_svm_6 = pickle.load(open('./lsa_svm_model/model_svm_6.pkl', 'rb'))
model_svm_7 = pickle.load(open('./lsa_svm_model/model_svm_7.pkl', 'rb'))
model_svm_8 = pickle.load(open('./lsa_svm_model/model_svm_8.pkl', 'rb'))
model_svm_9 = pickle.load(open('./lsa_svm_model/model_svm_9.pkl', 'rb'))
model_svm_10 = pickle.load(open('./lsa_svm_model/model_svm_10.pkl', 'rb'))

model_vectorizer_svm_1 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_1.pkl', 'rb'))
model_vectorizer_svm_2 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_2.pkl', 'rb'))
model_vectorizer_svm_3 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_3.pkl', 'rb'))
model_vectorizer_svm_4 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_4.pkl', 'rb'))
model_vectorizer_svm_5 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_5.pkl', 'rb'))
model_vectorizer_svm_6 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_6.pkl', 'rb'))
model_vectorizer_svm_7 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_7.pkl', 'rb'))
model_vectorizer_svm_8 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_8.pkl', 'rb'))
model_vectorizer_svm_9 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_9.pkl', 'rb'))
model_vectorizer_svm_10 = pickle.load(open('./lsa_svm_model/vectorizer_model_svm_10.pkl', 'rb'))

def process_similarity(message_data):
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

        (instansibaru, judulbaru) = text_finding_dokbaru.pdfparser(filename)
        new_document_title_list.append(judulbaru)

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
                text_proof = indikator09.ceklvl(filename)
            elif indicator_number == 10:
                text_proof = indikator10.ceklvl(filename)
            
            #HIGHLIGHT FOUND PROOF
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


            old_document_text_list = []
            old_document_title_list = []
            meeting_minutes_text_list = []
            meeting_minutes_title_list = []

            for old_document in old_document_list:
                    filename = "./src/static/" + old_document['name']
                    original_filename = old_document['original_name']
                    (instansilama, judullama) = text_finding_dokbaru.pdfparser(filename)
                    old_document_title_list.append(f'"{judullama}"')

            for idx, document in enumerate(old_document_list):
                
                filename = "./src/static/" + document['name']
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
                    text_proof = indikator09.ceklvl(filename)
                elif indicator_number == 10:
                    text_proof = indikator10.ceklvl(filename)
                if text_proof != "":
                    old_document_text_list.append(f'"{text_proof}"')
            
            for meeting_minutes in meeting_minutes_list:
                    filename = "./src/static/" + meeting_minutes['name']
                    original_filename = meeting_minutes['original_name']

                    #TODO notulen parser still needs to be improved to return both title and content
                    meeting_minutes_title = text_finding_notulen.pdfparser(filename)
                    meeting_minutes_title_list.append(f'"{meeting_minutes_title}"')
                    meeting_minutes_text_list.append(f'"{meeting_minutes_title}"')

            logging.warning("DATA PREDICTION")
            data_prediction = {
                'indikator': [indicator_number],
                'Judul': [judulbaru],
                'teks': [text_document_proof],
                'JudulDokumenLama': [f'[{",".join(old_document_title_list)}]'],
                'TeksDokumenLama': [f'[{",".join(old_document_text_list)}]'],
                'JudulNotulensiUndangan': [f'[{",".join(meeting_minutes_title_list)}]'],
                'IsiNotulensiUndangan': [f'[{",".join(meeting_minutes_text_list)}]']
            }
            logging.warning(data_prediction)
            df_predict = pd.DataFrame(data_prediction)
            for index, row in df_predict.iterrows():
                df_per_row = row.to_frame().transpose()
                if indicator_number == 1:
                    logging.warning("prediksi indikator 1")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_1 ,model_lsa_svm_1, model_svm_1)[0]
                elif indicator_number == 2:
                    logging.warning("prediksi indikator 2")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_2 ,model_lsa_svm_2, model_svm_2)[0]
                elif indicator_number == 3:
                    logging.warning("prediksi indikator 3")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_3 ,model_lsa_svm_3, model_svm_3)[0]
                elif indicator_number == 4:
                    logging.warning("prediksi indikator 4")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_4 ,model_lsa_svm_4, model_svm_4)[0]
                elif indicator_number == 5:
                    logging.warning("prediksi indikator 5")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_5 ,model_lsa_svm_5, model_svm_5)[0]
                elif indicator_number == 6:
                    logging.warning("prediksi indikator 6")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_6 ,model_lsa_svm_6, model_svm_6)[0]
                elif indicator_number == 7:
                    logging.warning("prediksi indikator 7")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_7 ,model_lsa_svm_7, model_svm_7)[0]
                elif indicator_number == 8:
                    logging.warning("prediksi indikator 8")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_8 ,model_lsa_svm_8, model_svm_8)[0]
                elif indicator_number == 9:
                    logging.warning("prediksi indikator 9")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_9 ,model_lsa_svm_9, model_svm_9)[0]
                elif indicator_number == 10:
                    logging.warning("prediksi indikator 10")
                    level = ts_preprocess.predict_model(df_per_row, model_vectorizer_svm_10 ,model_lsa_svm_10, model_svm_10)[0]

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
    
    payload = {
        "user_id": message_data['user_id'],
        "assessment_id": message_data['assessment_id'],
        "recipient_number": message_data['recipient_number'],
        "data": result_list
    }
    
    return payload