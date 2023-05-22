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

import ast
import math

import pickle
from sklearn import pipeline
from sklearn import linear_model
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.tree import DecisionTreeClassifier
from sklearn.metrics import classification_report
from sklearn.tree import export_graphviz
import graphviz

from Sastrawi.Stemmer.StemmerFactory import StemmerFactory

nltk.download('stopwords')
nltk.download('punkt')
nltk.download('averaged_perceptron_tagger')
nltk.download('wordnet')
nltk.download('omw-1.4')

stopword_list = stopwords.words('indonesian')
stopword_list.remove("setiap")
stopword_list.remove("seluruh")

factory = StemmerFactory()
stemmer = factory.create_stemmer()

def lowercase_sentence(text):
  if (type(text) == str):
    text = [text]
  
  for i in range(len(text)):
    text[i] = text[i].lower()
  
  return text

def convert_list(val):
    try:
        val = ast.literal_eval(val)
        if (val == ""):
            return None
        else:
            return val
    except (SyntaxError, ValueError):
        return val

def remove_stopwords(text):
  if (type(text) == str):
    text = [text]

  hasil = []
  for i in range(len(text)):
    tokens = word_tokenize(text[i])
  
    filtered_tokens = [token for token in tokens if token not in stopword_list]
    
    filtered_text = ' '.join(filtered_tokens)

    hasil.append(filtered_text)
  
  
  return hasil

def repair_type_for_model(df):
  df['indikator'] = df['indikator'].astype(int)
  df['Judul'] = df['Judul'].astype(str)
  df['teks'] = df['teks'].astype(str)
  df['JudulDokumenLama'] = df['JudulDokumenLama'].apply(convert_list)
  df['TeksDokumenLama'] = df['TeksDokumenLama'].apply(convert_list)
  df['JudulNotulensiUndangan'] = df['JudulNotulensiUndangan'].apply(convert_list)
  df['IsiNotulensiUndangan'] = df['IsiNotulensiUndangan'].apply(convert_list)

def lemmatization_sentence(text):
  if (type(text) == str):
    text = [text]

  hasil = []
  for i in range(len(text)):
    tokens = text[i].split()
    
    lemmatized_tokens = [stemmer.stem(word) for word in tokens]
    lemmatized_text = ' '.join(lemmatized_tokens)

    hasil.append(lemmatized_text)

  return hasil

def apply_preprocess_for_model(df):
  for col_name in df.columns[1:]:
    df[col_name]=df[col_name].apply(lowercase_sentence)
    df[col_name]=df[col_name].apply(remove_special_characters)
    df[col_name]=df[col_name].apply(remove_stopwords)
    df[col_name]=df[col_name].apply(lemmatization_sentence)

def remove_special_characters(text, remove_digits=True):
  if (type(text) == str):
    text = [text]

  for i in range(len(text)):
    pattern=r'[^a-zA-z0-9\s]'
    text[i]=re.sub(pattern,'',text[i])
  
  return text

def predict_model(df_predict, vectorizer, lsa, svm):
  repair_type_for_model(df_predict)
  apply_preprocess_for_model(df_predict)

  teks = df_predict['teks'].apply(lambda x: x[0])
  isiNotulensiUndangan = df_predict['IsiNotulensiUndangan'].apply(lambda x: ' '.join(x))

  tfidf_predict = vectorizer.transform(teks + isiNotulensiUndangan)

  lsa_predict = lsa.transform(tfidf_predict)

  predicted_level = svm.predict(lsa_predict)

  print("Prediksi level: ", predicted_level)

  return predicted_level