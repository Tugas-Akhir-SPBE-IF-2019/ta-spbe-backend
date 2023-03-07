# TF
# this code is for TEXT FINDING - INDIKATOR 1
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
import preprocess_dokbaru as dokbaru
import preprocess_doklama as doklama

#TEST TF MODEL 
filename = 'F2201-287-Indikator_01_+_Indikator1_Perbup_81_tahun_2021.pdf'
filename = 'Draft Perbup SPBE 2021 Revisi.pdf'


def txtreader(filename, keyword):
    file = open(f'cleaned_{filename}.txt', 'r')

    idx = 0
    result = []

    lvl2_exclude = ['audit']

    for line in file:
        if (len(keyword) == 1):
            if re.search(keyword[0], line, re.IGNORECASE):
                result.append([idx, line])

        else:
            # check with keyword 'domain'
            for key in keyword:
                reg = f'(?:(domain)\s(arsitektur)?\s({key}))'
                if re.search(reg, line, re.IGNORECASE):
                    result.append([idx, line])

            # TODO
            # check without keyword 'domain'

        idx += 1

    file.close()
    return (result)


# if __name__ == '__main__':
#     pdfparser(sys.argv[1])

def ceklvl():
    list_final = []
    text_final = ''

    lvl1 = ["arsitektur SPBE"]
    res1 = txtreader(filename, lvl1)

    # cek if keyword lvl1 is not found
    if (not res1):
        return ''

    lvl2 = ["Proses Bisnis", "Data dan Informasi", "Infrastruktur SPBE",
            "Aplikasi SPBE", "Keamanan SPBE", "Layanan SPBE"]
    res2 = txtreader(filename, lvl2)

    if (not res2):
        return ''

    for el in res2:
        list_final.append(el[1])

    lvl4 = ["integrasi", "reviu"]
    res4 = txtreader(filename, lvl4)

    for el in res4:
        list_final.append(el[1])

    text_final = ". ".join(list_final)

    # clean text
    text_final = re.sub(r'(\n)+', '', text_final, flags=re.MULTILINE)
    text_final = re.sub(r'(;)+', ',', text_final, flags=re.MULTILINE)

    return text_final


# filename = 'F2201-287-Indikator_01_+_Indikator1_Perbup_81_tahun_2021.pdf'
# filename = 'Draft Perbup SPBE 2021 Revisi.pdf'

# (instansibaru, judulbaru) = dokbaru.pdfparser(filename)
# (instansilama, judullama) = doklama.pdfparser(filename)

# level_final = ceklvl()
# print(level_final)

# TF

import nsq, ast, toml, logging, requests, threading

# Import libraries
import json
import pandas as pd
import numpy as np

import matplotlib.pyplot as plt
import seaborn as sns

from sklearn import metrics, manifold
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report, confusion_matrix,accuracy_score
from sklearn.feature_extraction.text import TfidfVectorizer

from sklearn.feature_extraction.text import TfidfTransformer
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.metrics.pairwise import cosine_similarity

import re

import transformers

import nltk
from nltk.tokenize import word_tokenize
from nltk.corpus import stopwords
from nltk.stem import WordNetLemmatizer
# Stemmer nltk
from nltk.stem import PorterStemmer

# Stemmer sastrawi
from Sastrawi.Stemmer.StemmerFactory import StemmerFactory

import pickle

model = pickle.load(open('model.pckl', 'rb'))


def send_result(message_data, config):
    #TEST TF MODEL

    (instansibaru, judulbaru) = dokbaru.pdfparser(filename)
    (instansilama, judullama) = doklama.pdfparser(filename)

    level_final = ceklvl()

    pred_result = model.predict([level_final])    
    #TEST TF MODEL

    callback_endpoint = 'http://' + config['server']['host'] + '/assessments/result/callback'
    payload = {
        "user_id": message_data['UserId'],
        "assessment_id": message_data['AssessmentId'],
        "recipient_number": message_data['RecipientNumber'],
        "indicator_assessment_id": message_data['IndicatorAssessmentId'],
        "level": int(pred_result[0]),
        "explanation": "berdasarkan data dukung yang diberikan, level yang sesuai adalah level " + str(pred_result[0]),
        "support_data_document_id": message_data['Content'],
        "proof": level_final
        
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

# nltk.download('stopwords')
# nltk.download('punkt')
# nltk.download('averaged_perceptron_tagger')
# nltk.download('wordnet')
# nltk.download('omw-1.4')


# stopword_list = stopwords.words('indonesian')
# lemmatizer = WordNetLemmatizer()
# stemmer = PorterStemmer()

# factory = StemmerFactory()
# stemmer = factory.create_stemmer()

# def lowercase_sentence(text):
#     text = text.lower()

#     return text

# def remove_between_square_brackets(text):
#     return re.sub('\[[^]]*\]', '', text)

# def remove_special_characters(text, remove_digits=True):
#     pattern=r'[^a-zA-z0-9\s]'
#     text=re.sub(pattern,'',text)
#     return text

# def remove_stopwords(text, is_lower_case=False):
#     tokens = word_tokenize(text)

#     if is_lower_case:
#         filtered_tokens = [token for token in tokens if token not in stopword_list]
#     else:
#         filtered_tokens = [token for token in tokens if token.lower() not in stopword_list]

#     filtered_text = ' '.join(filtered_tokens)
#     return filtered_text

# def lemmatization_sentence(text):
#     tokens = word_tokenize(text)

#     lemmatize_tokens = [lemmatizer.lemmatize(token) for token in tokens if not token in set(stopword_list)]
#     lemmatize_text = ' '.join(lemmatize_tokens)

#     # print(tokens, lemmatize_tokens, lemmatize_text)

#     return text

# def stemming_sentence(text):
#     tokens = word_tokenize(text)

#     stemming_text = stemmer.stem(text)

#     return stemming_text

# def apply_preprocess_to_model(list_teks):
#     hasil = []
#     for teks in list_teks:
#         teks = lowercase_sentence(teks)
#         teks = remove_between_square_brackets(teks)
#         teks = remove_special_characters(teks)
#         teks = remove_stopwords(teks)
#         teks = lemmatization_sentence(teks)
#         teks = stemming_sentence(teks)

#         hasil.append(teks)

#     return hasil

# teks = [
#     'Arsitektur SPBE bertujuan untuk memandu implementasi integrasi proses Bisnis, data dan informasi, Arsitektur SPBE, aplikasi SPBE, dan Keamanan SPBE untuk menghadirkan Layanan SPBE terintegrasi. Arsitektur SPBE Pemerintah Daerah meliputi: domain Arsitektur Proses Bisnis, data dan informasi tentang Domain Infrastruktur, domain Infrastruktur SPBE, domain Arsitektur Aplikasi SPBE, domain Arsitektur Pengamanan SPBE, dan domain Arsitektur Layanan SPBE. Kerangka SPBE Pemerintah Daerah akan ditinjau ulang pada jangka menengah dan pada tahun terakhir pelaksanaan, atau lebih sering sesuai kebutuhan. Rencana Pembangunan Daerah. Tinjauan Struktur SPBE Pemerintah Daerah sebagaimana dimaksud pada ayat (1) disusun oleh Sekretaris Daerah.',
#     'Tujuan dari arsitektur SPBE adalah untuk memberikan panduan integrasi proses bisnis, data dan informasi, infrastruktur SPBE, aplikasi SPBE, dan keamanan SPBE untuk menghasilkan layanan terintegrasi SPBE. Arsitektur SPBE pemerintah daerah meliputi: domain arsitektur proses bisnis, domain arsitektur data dan informasi, domain arsitektur infrastruktur SPBE, domain arsitektur aplikasi SPBE, domain arsitektur keamanan SPBE, dan domain arsitektur layanan SPBE. Susunan SPBE Pemda ditinjau jangka menengah dan pada tahun terakhir pelaksanaan atau sewaktu-waktu sesuai kebutuhan. Rencana Pembangunan Jangka Menengah Daerah. Tinjauan struktur SPBE pemerintah daerah sebagaimana dimaksud pada ayat (1) dikoordinasikan oleh Sekretaris Daerah.',
#     'Arsitektur SPBE dimaksudkan untuk memberikan panduan bagaimana mengimplementasikan integrasi proses bisnis, data dan informasi, infrastruktur SPBE, aplikasi SPBE dan keamanan SPBE untuk menghasilkan layanan SPBE yang terintegrasi. SPBE Pemda meliputi: Domain Arsitektur Proses Bisnis, Domain Arsitektur Data dan Informasi, Domain Arsitektur Infrastruktur SPBE, Domain Arsitektur Aplikasi SPBE, Domain Arsitektur Keamanan SPBE, dan Domain Arsitektur Layanan SPBE. Peninjauan Arsitektur SPBE Pemerintahan Mandiri Voivodeship dilakukan paruh waktu dan pada tahun terakhir pelaksanaan atau sewaktu-waktu tergantung kebutuhan Unsur-unsur SPBE sebagaimana dimaksud dalam Pasal. Rencana Pembangunan Jangka Menengah Daerah. Tinjauan Arsitektur SPBE dari Pemerintahan Sendiri Voivodeship, sebagaimana dimaksud pada par. 1 dikoordinir oleh Sekretaris Voivodship.'
# ]
# teks = apply_preprocess_to_model(teks)
# print(model.predict(teks))

# TEST

r = reader = nsq.Reader(
            topic=nsq_topic, channel=nsq_channel, message_handler=message_handler,
            lookupd_connect_timeout=10, requeue_delay=10, 
            nsqd_tcp_addresses=[nsqd_address], max_in_flight=5, snappy=False
    )
nsq.run()