import unittest
import app.text_similarity.process_similarity as similarity

class TestApp(unittest.TestCase):
    def test_process_similarity_all_document_provided_valid(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                # {
                #     'id': '77260ad0-3c0a-441a-8dec-783c2293d507', 
                #     'number': 2, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Peta Rencana SPBE Instansi Pusat/Pemerintah Daerah'
                # }, 
                # {
                #     'id': '0d6637fd-7cb7-498d-a34d-e74dc4f4bde1', 
                #     'number': 3, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Manajemen Data di Instansi Pusat/Pemerintah Daerah'
                # }, 
                # {
                #     'id': '0cf24829-3bed-4995-9639-7ff8ca4166ba', 
                #     'number': 4, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Pembangunan Aplikasi SPBE'
                # }, 
                # {
                #     'id': '1691d04b-45c5-4e18-a6f7-52bae90b7182', 
                #     'number': 5, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Layanan Pusat Data'
                # }, 
                # {
                #     'id': 'aff23ec0-9777-4246-8c20-1a36c0b95162', 
                #     'number': 6, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Layanan Jaringan Intra Instansi Pusat/Pemerintah Daerah'
                # }, 
                # {
                #     'id': 'b7b908b9-d215-411f-86b5-94b9ab660b64', 
                #     'number': 7, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Penggunaan Sistem Penghubung Layanan Instansi Pusat/Pemerintah Daerah'
                # }, 
                # {
                #     'id': 'c76978e8-7f24-4b1e-9574-b4cab882ed13', 
                #     'number': 8, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Manajemen Keamanan Informasi. Evidence kebijakan internal terkait Manajemen Keamanan Informasi'
                # }, 
                # {
                #     'id': 'e67f5be3-1797-4af9-b338-3733a8e8c3e9', 
                #     'number': 9, 
                #     'detail': 'Kebijakan internal Audit TIK'
                # }, 
                # {
                #     'id': 'ffec2af3-f311-403c-8138-63e383c84305', 
                #     'number': 10, 
                #     'detail': 'Tingkat Kematangan Kebijakan Internal Tim Koordinasi SPBE Instansi Pusat/Pemerintah Daerah'
                # }
                ], 
            'support_document_list': [
                {
                    'name': 'test_new_document_valid.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                }, 
                {
                    'name': 'test_old_document_valid.pdf', 
                    'original_name': 'BatuBara.pdf', 
                    'type': 'OLD_DOCUMENT'
                }, 
                {
                    'name': 'test_meeting_minutes_valid.pdf', 
                    'original_name': 'Nota Dinas.pdf', 
                    'type': 'MEETING_MINUTES'
                }], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])

    def test_process_similarity_single_new_document_provided_valid(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                ], 
            'support_document_list': [
                {
                    'name': 'test_new_document_valid.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                }, 
                # {
                #     'name': 'test_old_document_valid.pdf', 
                #     'original_name': 'BatuBara.pdf', 
                #     'type': 'OLD_DOCUMENT'
                # }, 
                # {
                #     'name': 'test_meeting_minutes_valid.pdf', 
                #     'original_name': 'Nota Dinas.pdf', 
                #     'type': 'MEETING_MINUTES'
                # }
                ], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])
    def test_process_similarity_single_new_and_old_document_provided_valid(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                ], 
            'support_document_list': [
                {
                    'name': 'test_new_document_valid.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                }, 
                {
                    'name': 'test_old_document_valid.pdf', 
                    'original_name': 'BatuBara.pdf', 
                    'type': 'OLD_DOCUMENT'
                }, 
                ], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])

    def test_process_similarity_single_new_and_meeting_document_provided_valid(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                ], 
            'support_document_list': [
                {
                    'name': 'test_new_document_valid.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                }, 
                {
                    'name': 'test_meeting_minutes_valid.pdf', 
                    'original_name': 'Nota Dinas.pdf', 
                    'type': 'MEETING_MINUTES'
                }
                ], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])

    def test_process_similarity_single_new_document_provided_non_spbe(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                ], 
            'support_document_list': [
                {
                    'name': 'test_non_spbe_document.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                }, 
                ], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])

    def test_process_similarity_single_new_document_spbe_and_old_document_non_spbe_provided(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                ], 
            'support_document_list': [
                {
                    'name': 'test_new_document_valid.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                },
                {
                    'name': 'test_non_spbe_document.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'OLD_DOCUMENT'
                }, 
                ], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])

    def test_process_similarity_single_new_document_spbe_and_meeting_minutes_non_spbe_provided(self):
        mock_message_data = {
            'name': 'Kabupaten Batu Bara', 
            'content': 'c75a3298-ab99-4aa2-960f-778c4eba0acf', 
            'user_id': 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 
            'recipient_number': '6285157017311', 
            'assessment_id': 'e16db86b-4ed1-46cd-9065-dc4ba3a7774f', 
            'indicator_assessment_list': [
                {
                    'id': '1b0b9a3b-25a9-4a70-9cca-286f0e18652d', 
                    'number': 1, 
                    'detail': 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah'
                }, 
                ], 
            'support_document_list': [
                {
                    'name': 'test_new_document_valid.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'NEW_DOCUMENT'
                },
                {
                    'name': 'test_non_spbe_document.pdf', 
                    'original_name': 'Draft Perbup SPBE 2021 Revisi.pdf', 
                    'type': 'MEETING_MINUTES'
                }, 
                ], 
            'indicator_detail': '', 
            'institution_name': 'Kabupaten Batu Bara', 
            'timestamp': '2023-05-22 12:57:53.034203843 +0000 UTC'}
        similarity_result = similarity.process_similarity(mock_message_data)

        #assert that prediction result is not 0
        self.assertNotEqual(0, similarity_result["data"][0]["result"]["level"])