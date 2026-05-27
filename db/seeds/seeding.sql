TRUNCATE TABLE
    topup_details, transfer_details, transactions,
    favorite_contacts, newsletter, reviews,
    forgot_password, oauth_user, wallet, users,
    payment_method, category_payment_method
RESTART IDENTITY CASCADE;

-- USERS
INSERT INTO users (id, full_name, email, hash_password, hash_pin, phone, profile_picture_url, is_verified) VALUES
                                                                                                               (1,  'Budi Santoso',  'budi.santoso@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281111111111', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Budi', TRUE),
                                                                                                               (2,  'Siti Aminah',   'siti.aminah@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281222222222', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Siti', TRUE),
                                                                                                               (3,  'Andi Wijaya',   'andi.wijaya@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281333333333', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Andi', TRUE),
                                                                                                               (4,  'Dewi Lestari',  'dewi.lestari@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281444444444', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Dewi', FALSE),
                                                                                                               (5,  'Eko Prasetyo',  'eko.prasetyo@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281555555555', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Eko', TRUE),
                                                                                                               (6,  'Fajar Hidayat', 'fajar.hidayat@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281666666006', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Fajar', TRUE),
                                                                                                               (7,  'Gita Kusuma',   'gita.kusuma@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281666666007', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Gita', TRUE),
                                                                                                               (8,  'Hendra Putra',  'hendra.putra@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281666666008', NULL, FALSE),
                                                                                                               (9,  'Indah Pratama', 'indah.pratama@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281666666009', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Indah', TRUE),
                                                                                                               (10, 'Joko Santoso',  'joko.santoso@email.com',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$xIJYsW0jYK6HncD8JVUsZA$wGP3m49Qm7OrrFOAu/baOcutE3s/sUb0LaoGLy7B13Y',
                                                                                                                '$argon2id$v=19$m=65536,t=2,p=1$I5cioGD7oLSfhvYC6NfroQ$LBPR9U069Fz//S9r70hTalg57q340Nnsn2UgKpqPVjI',
                                                                                                                '+6281666666010', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Joko', TRUE);

-- WALLET
INSERT INTO wallet (user_id, balance, updated_at) VALUES
                                                      (1,  10000000.00, NOW()),
                                                      (2,   2000000.00, NOW()),
                                                      (3,   1500000.00, NOW()),
                                                      (4,   1000000.00, NOW()),
                                                      (5,     50000.00, NOW()),
                                                      (6,    800000.00, NOW()),
                                                      (7,    600000.00, NOW()),
                                                      (8,    200000.00, NOW()),
                                                      (9,    500000.00, NOW()),
                                                      (10,   700000.00, NOW());

-- CATEGORY PAYMENT METHOD
INSERT INTO category_payment_method (id, category_name) VALUES
                                                            (1, 'Bank Transfer'),
                                                            (2, 'E-Wallet'),
                                                            (3, 'Convenience Store'),
                                                            (4, 'Credit Card'),
                                                            (5, 'Instant Payment');

-- PAYMENT METHOD
INSERT INTO payment_method (id, payment_category_id, name, code, fee, logo_url, is_active) VALUES
                                                                                               (1,  1, 'BCA Virtual Account',     'BCAVA',    1000, 'https://logo.url/bcava.png',    TRUE),
                                                                                               (2,  1, 'Mandiri Virtual Account', 'MANDIRIVA',1000, 'https://logo.url/mandiriva.png',TRUE),
                                                                                               (3,  1, 'BNI Virtual Account',     'BNIVA',    1000, 'https://logo.url/bniva.png',    TRUE),
                                                                                               (4,  1, 'BRI Virtual Account',     'BRIVA',    1000, 'https://logo.url/briva.png',    TRUE),
                                                                                               (5,  2, 'GoPay',                   'GOPAY',    1500, 'https://logo.url/gopay.png',    TRUE),
                                                                                               (6,  2, 'OVO',                     'OVO',      1500, 'https://logo.url/ovo.png',      TRUE),
                                                                                               (7,  2, 'Dana',                    'DANA',     1000, 'https://logo.url/dana.png',     TRUE),
                                                                                               (8,  3, 'Indomaret',               'IDM',      2500, 'https://logo.url/idm.png',      TRUE),
                                                                                               (9,  3, 'Alfamart',                'ALFA',     2500, 'https://logo.url/alfa.png',     TRUE),
                                                                                               (10, 4, 'Visa Card',               'VISA',     5000, 'https://logo.url/visa.png',     TRUE),
                                                                                               (11, 5, 'QRIS',                    'QRIS',        0, 'https://logo.url/qris.png',     TRUE),
                                                                                               (12, 5, 'BI-Fast',                 'BIFAST',   2500, 'https://logo.url/bifast.png',   TRUE);

-- ============================================================
-- TRANSACTIONS
-- Konvensi:
--   activity_type='topup'    → type='income', TIDAK masuk report/history transfer
--   activity_type='transfer' → type='expense' (sender) atau type='income' (receiver)
--
-- Query GetTransactionReport dan GetTransactionHistory
-- WAJIB filter: WHERE activity_type = 'transfer'
-- ============================================================

-- TOPUP (diabaikan di report & history, hanya untuk balance context)
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (1,  1, 5000000, 'income', 'topup', 'success', NOW() - INTERVAL '400 days'),
                                                                                            (2,  1, 3000000, 'income', 'topup', 'success', NOW() - INTERVAL '200 days'),
                                                                                            (3,  1, 5000000, 'income', 'topup', 'success', NOW() - INTERVAL '60 days'),
                                                                                            (4,  2, 3000000, 'income', 'topup', 'success', NOW() - INTERVAL '50 days'),
                                                                                            (5,  3, 2000000, 'income', 'topup', 'success', NOW() - INTERVAL '45 days'),
                                                                                            (6,  5,   50000, 'income', 'topup', 'success', NOW() - INTERVAL '5 days');

INSERT INTO topup_details (transaction_id, payment_method_id, order_amount, delivery_fee, tax_amount, total_amount) VALUES
                                                                                                                        (1, 1,  5000000, 1000, 0, 5001000),
                                                                                                                        (2, 11, 3000000,    0, 0, 3000000),
                                                                                                                        (3, 1,  5000000, 1000, 0, 5001000),
                                                                                                                        (4, 5,  3000000, 1500, 0, 3001500),
                                                                                                                        (5, 7,  2000000, 1000, 0, 2001000),
                                                                                                                        (6, 11,   50000,    0, 0,   50000);

-- ============================================================
-- TRANSFER — 12 BULAN LALU (untuk chart period=year)
-- Budi (1) sebagai pusat, berinteraksi dengan user 2,3,6,7,9,10
-- ============================================================

-- Bulan 12 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (7,  1, 800000, 'expense', 'transfer', 'success', NOW() - INTERVAL '365 days'),
                                                                                            (8,  2, 800000, 'income',  'transfer', 'success', NOW() - INTERVAL '365 days'),
                                                                                            (9,  3, 500000, 'expense', 'transfer', 'success', NOW() - INTERVAL '360 days'),
                                                                                            (10, 1, 500000, 'income',  'transfer', 'success', NOW() - INTERVAL '360 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (7,  2, 'Transfer bulan 12 lalu'),(8,  1, 'Transfer bulan 12 lalu'),
                                                                            (9,  1, 'Bayar hutang bulan 12'), (10, 3, 'Bayar hutang bulan 12');

-- Bulan 11 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (11, 1, 1200000, 'expense', 'transfer', 'success', NOW() - INTERVAL '335 days'),
                                                                                            (12, 2, 1200000, 'income',  'transfer', 'success', NOW() - INTERVAL '335 days'),
                                                                                            (13, 6, 300000,  'expense', 'transfer', 'success', NOW() - INTERVAL '330 days'),
                                                                                            (14, 1, 300000,  'income',  'transfer', 'success', NOW() - INTERVAL '330 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (11, 2, 'Transfer bulan 11 lalu'),(12, 1, 'Transfer bulan 11 lalu'),
                                                                            (13, 1, 'Kembalian bulan 11'),    (14, 6, 'Kembalian bulan 11');

-- Bulan 10 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (15, 1, 900000, 'expense', 'transfer', 'success', NOW() - INTERVAL '305 days'),
                                                                                            (16, 3, 900000, 'income',  'transfer', 'success', NOW() - INTERVAL '305 days'),
                                                                                            (17, 2, 400000, 'expense', 'transfer', 'success', NOW() - INTERVAL '300 days'),
                                                                                            (18, 1, 400000, 'income',  'transfer', 'success', NOW() - INTERVAL '300 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (15, 3, 'Transfer bulan 10 lalu'),(16, 1, 'Transfer bulan 10 lalu'),
                                                                            (17, 1, 'Kembalian bulan 10'),    (18, 2, 'Kembalian bulan 10');

-- Bulan 9 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (19, 1, 1500000, 'expense', 'transfer', 'success', NOW() - INTERVAL '275 days'),
                                                                                            (20, 2, 1500000, 'income',  'transfer', 'success', NOW() - INTERVAL '275 days'),
                                                                                            (21, 7, 200000,  'expense', 'transfer', 'success', NOW() - INTERVAL '270 days'),
                                                                                            (22, 1, 200000,  'income',  'transfer', 'success', NOW() - INTERVAL '270 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (19, 2, 'Transfer bulan 9 lalu'), (20, 1, 'Transfer bulan 9 lalu'),
                                                                            (21, 1, 'Kembalian bulan 9'),     (22, 7, 'Kembalian bulan 9');

-- Bulan 8 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (23, 1, 700000, 'expense', 'transfer', 'success', NOW() - INTERVAL '245 days'),
                                                                                            (24, 6, 700000, 'income',  'transfer', 'success', NOW() - INTERVAL '245 days'),
                                                                                            (25, 3, 350000, 'expense', 'transfer', 'success', NOW() - INTERVAL '240 days'),
                                                                                            (26, 1, 350000, 'income',  'transfer', 'success', NOW() - INTERVAL '240 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (23, 6, 'Transfer bulan 8 lalu'), (24, 1, 'Transfer bulan 8 lalu'),
                                                                            (25, 1, 'Kembalian bulan 8'),     (26, 3, 'Kembalian bulan 8');

-- Bulan 7 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (27, 1, 2000000, 'expense', 'transfer', 'success', NOW() - INTERVAL '215 days'),
                                                                                            (28, 2, 2000000, 'income',  'transfer', 'success', NOW() - INTERVAL '215 days'),
                                                                                            (29, 9, 500000,  'expense', 'transfer', 'success', NOW() - INTERVAL '210 days'),
                                                                                            (30, 1, 500000,  'income',  'transfer', 'success', NOW() - INTERVAL '210 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (27, 2, 'Transfer bulan 7 lalu'), (28, 1, 'Transfer bulan 7 lalu'),
                                                                            (29, 1, 'Kembalian bulan 7'),     (30, 9, 'Kembalian bulan 7');

-- Bulan 6 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (31, 1, 1100000, 'expense', 'transfer', 'success', NOW() - INTERVAL '185 days'),
                                                                                            (32, 3, 1100000, 'income',  'transfer', 'success', NOW() - INTERVAL '185 days'),
                                                                                            (33, 6, 250000,  'expense', 'transfer', 'success', NOW() - INTERVAL '180 days'),
                                                                                            (34, 1, 250000,  'income',  'transfer', 'success', NOW() - INTERVAL '180 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (31, 3, 'Transfer bulan 6 lalu'), (32, 1, 'Transfer bulan 6 lalu'),
                                                                            (33, 1, 'Kembalian bulan 6'),     (34, 6, 'Kembalian bulan 6');

-- Bulan 5 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (35, 1, 800000, 'expense', 'transfer', 'success', NOW() - INTERVAL '155 days'),
                                                                                            (36, 7, 800000, 'income',  'transfer', 'success', NOW() - INTERVAL '155 days'),
                                                                                            (37, 2, 600000, 'expense', 'transfer', 'success', NOW() - INTERVAL '150 days'),
                                                                                            (38, 1, 600000, 'income',  'transfer', 'success', NOW() - INTERVAL '150 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (35, 7, 'Transfer bulan 5 lalu'), (36, 1, 'Transfer bulan 5 lalu'),
                                                                            (37, 1, 'Kembalian bulan 5'),     (38, 2, 'Kembalian bulan 5');

-- Bulan 4 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (39, 1, 1800000, 'expense', 'transfer', 'success', NOW() - INTERVAL '125 days'),
                                                                                            (40, 2, 1800000, 'income',  'transfer', 'success', NOW() - INTERVAL '125 days'),
                                                                                            (41, 3, 400000,  'expense', 'transfer', 'success', NOW() - INTERVAL '120 days'),
                                                                                            (42, 1, 400000,  'income',  'transfer', 'success', NOW() - INTERVAL '120 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (39, 2, 'Transfer bulan 4 lalu'), (40, 1, 'Transfer bulan 4 lalu'),
                                                                            (41, 1, 'Kembalian bulan 4'),     (42, 3, 'Kembalian bulan 4');

-- Bulan 3 lalu
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (43, 1,  950000, 'expense', 'transfer', 'success', NOW() - INTERVAL '95 days'),
                                                                                            (44, 6,  950000, 'income',  'transfer', 'success', NOW() - INTERVAL '95 days'),
                                                                                            (45, 9,  300000, 'expense', 'transfer', 'success', NOW() - INTERVAL '90 days'),
                                                                                            (46, 1,  300000, 'income',  'transfer', 'success', NOW() - INTERVAL '90 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (43, 6, 'Transfer bulan 3 lalu'), (44, 1, 'Transfer bulan 3 lalu'),
                                                                            (45, 1, 'Kembalian bulan 3'),     (46, 9, 'Kembalian bulan 3');

-- ============================================================
-- TRANSFER — BULAN LALU (untuk chart period=month, minggu 1-4)
-- ============================================================

-- Minggu 4 bulan lalu (sekitar 28-35 hari lalu)
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (47, 1, 500000, 'expense', 'transfer', 'success', NOW() - INTERVAL '35 days'),
                                                                                            (48, 2, 500000, 'income',  'transfer', 'success', NOW() - INTERVAL '35 days'),
                                                                                            (49, 3, 200000, 'expense', 'transfer', 'success', NOW() - INTERVAL '32 days'),
                                                                                            (50, 1, 200000, 'income',  'transfer', 'success', NOW() - INTERVAL '32 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (47, 2, 'Patungan minggu 4'),(48, 1, 'Patungan minggu 4'),
                                                                            (49, 1, 'Kembalian minggu 4'),(50, 3, 'Kembalian minggu 4');

-- Minggu 3 bulan lalu (sekitar 21-28 hari lalu)
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (51, 1, 750000, 'expense', 'transfer', 'success', NOW() - INTERVAL '27 days'),
                                                                                            (52, 6, 750000, 'income',  'transfer', 'success', NOW() - INTERVAL '27 days'),
                                                                                            (53, 7, 150000, 'expense', 'transfer', 'success', NOW() - INTERVAL '24 days'),
                                                                                            (54, 1, 150000, 'income',  'transfer', 'success', NOW() - INTERVAL '24 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (51, 6, 'Patungan minggu 3'),(52, 1, 'Patungan minggu 3'),
                                                                            (53, 1, 'Kembalian minggu 3'),(54, 7, 'Kembalian minggu 3');

-- Minggu 2 bulan lalu (sekitar 14-21 hari lalu)
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (55, 1,  600000, 'expense', 'transfer', 'success', NOW() - INTERVAL '20 days'),
                                                                                            (56, 2,  600000, 'income',  'transfer', 'success', NOW() - INTERVAL '20 days'),
                                                                                            (57, 9,  250000, 'expense', 'transfer', 'success', NOW() - INTERVAL '17 days'),
                                                                                            (58, 1,  250000, 'income',  'transfer', 'success', NOW() - INTERVAL '17 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (55, 2, 'Patungan minggu 2'),(56, 1, 'Patungan minggu 2'),
                                                                            (57, 1, 'Kembalian minggu 2'),(58, 9, 'Kembalian minggu 2');

-- Minggu 1 bulan lalu (sekitar 7-14 hari lalu)
INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (59, 1, 1000000, 'expense', 'transfer', 'success', NOW() - INTERVAL '13 days'),
                                                                                            (60, 3, 1000000, 'income',  'transfer', 'success', NOW() - INTERVAL '13 days'),
                                                                                            (61, 6,  300000, 'expense', 'transfer', 'success', NOW() - INTERVAL '10 days'),
                                                                                            (62, 1,  300000, 'income',  'transfer', 'success', NOW() - INTERVAL '10 days');
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (59, 3, 'Patungan minggu 1'),(60, 1, 'Patungan minggu 1'),
                                                                            (61, 1, 'Kembalian minggu 1'),(62, 6, 'Kembalian minggu 1');

-- ============================================================
-- TRANSFER — MINGGU INI (untuk chart period=week, hari senin-minggu)
-- ============================================================

INSERT INTO transactions (id, user_id, amount, type, activity_type, status, created_at) VALUES
                                                                                            (63, 1, 200000, 'expense', 'transfer', 'success', NOW() - INTERVAL '6 days'),
                                                                                            (64, 2, 200000, 'income',  'transfer', 'success', NOW() - INTERVAL '6 days'),
                                                                                            (65, 3, 100000, 'expense', 'transfer', 'success', NOW() - INTERVAL '5 days'),
                                                                                            (66, 1, 100000, 'income',  'transfer', 'success', NOW() - INTERVAL '5 days'),
                                                                                            (67, 1, 350000, 'expense', 'transfer', 'success', NOW() - INTERVAL '4 days'),
                                                                                            (68, 6, 350000, 'income',  'transfer', 'success', NOW() - INTERVAL '4 days'),
                                                                                            (69, 7,  75000, 'expense', 'transfer', 'success', NOW() - INTERVAL '3 days'),
                                                                                            (70, 1,  75000, 'income',  'transfer', 'success', NOW() - INTERVAL '3 days'),
                                                                                            (71, 1, 500000, 'expense', 'transfer', 'success', NOW() - INTERVAL '2 days'),
                                                                                            (72, 2, 500000, 'income',  'transfer', 'success', NOW() - INTERVAL '2 days'),
                                                                                            (73, 9, 200000, 'expense', 'transfer', 'success', NOW() - INTERVAL '1 days'),
                                                                                            (74, 1, 200000, 'income',  'transfer', 'success', NOW() - INTERVAL '1 days'),
                                                                                            (75, 1, 400000, 'expense', 'transfer', 'success', NOW()),
                                                                                            (76, 3, 400000, 'income',  'transfer', 'success', NOW());
INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES
                                                                            (63, 2, 'Senin - bayar makan'),   (64, 1, 'Senin - bayar makan'),
                                                                            (65, 1, 'Selasa - kembalian'),    (66, 3, 'Selasa - kembalian'),
                                                                            (67, 6, 'Rabu - ongkir'),         (68, 1, 'Rabu - ongkir'),
                                                                            (69, 1, 'Kamis - kembalian'),     (70, 7, 'Kamis - kembalian'),
                                                                            (71, 2, 'Jumat - transfer rutin'),(72, 1, 'Jumat - transfer rutin'),
                                                                            (73, 1, 'Sabtu - kembalian'),     (74, 9, 'Sabtu - kembalian'),
                                                                            (75, 3, 'Minggu - patungan'),     (76, 1, 'Minggu - patungan');

-- FAVORITE CONTACTS
INSERT INTO favorite_contacts (user_id, favorite_user_id) VALUES
                                                              (1, 2),(1, 3),(1, 6),(1, 7),
                                                              (2, 1),(2, 3),
                                                              (3, 1),(3, 2);

-- REVIEWS
INSERT INTO reviews (user_id, rating, description, created_at) VALUES
                                                                   (1, 5, 'Aplikasi sangat mudah digunakan!',           NOW() - INTERVAL '10 days'),
                                                                   (2, 4, 'Transfer lancar, tampilan bersih',           NOW() - INTERVAL '8 days'),
                                                                   (3, 5, 'Fitur riwayat transaksi lengkap',            NOW() - INTERVAL '6 days'),
                                                                   (6, 3, 'Bagus tapi kadang loading lambat',           NOW() - INTERVAL '5 days'),
                                                                   (7, 4, 'Mudah digunakan sehari-hari',                NOW() - INTERVAL '4 days'),
                                                                   (9, 5, 'Terbaik dibanding e-wallet lain',            NOW() - INTERVAL '3 days'),
                                                                   (10, 4, 'Saldo selalu aman dan akurat',              NOW() - INTERVAL '2 days'),
                                                                   (4, 2, 'Perlu peningkatan notifikasi transaksi',     NOW() - INTERVAL '1 days');

-- NEWSLETTER
INSERT INTO newsletter (user_id, email, status, created_at) VALUES
                                                                (1,  'budi.news@example.com',  'active',      NOW() - INTERVAL '30 days'),
                                                                (2,  'siti.news@example.com',  'active',      NOW() - INTERVAL '25 days'),
                                                                (3,  'andi.news@example.com',  'active',      NOW() - INTERVAL '20 days'),
                                                                (4,  'dewi.news@example.com',  'unsubscribe', NOW() - INTERVAL '15 days'),
                                                                (5,  'eko.news@example.com',   'active',      NOW() - INTERVAL '10 days'),
                                                                (6,  'fajar.news@example.com', 'active',      NOW() - INTERVAL '8 days'),
                                                                (7,  'gita.news@example.com',  'unsubscribe', NOW() - INTERVAL '5 days'),
                                                                (9,  'indah.news@example.com', 'active',      NOW() - INTERVAL '3 days'),
                                                                (10, 'joko.news@example.com',  'active',      NOW() - INTERVAL '1 days');

-- FORGOT PASSWORD
INSERT INTO forgot_password (user_id, token, is_used, created_at, expired_at) VALUES
                                                                                  (1, 'valid-token-budi-001',   FALSE, NOW() - INTERVAL '10 minutes', NOW() + INTERVAL '50 minutes'),
                                                                                  (2, 'expired-token-siti-002', FALSE, NOW() - INTERVAL '2 hours',    NOW() - INTERVAL '1 hours'),
                                                                                  (3, 'used-token-andi-003',    TRUE,  NOW() - INTERVAL '1 days',     NOW() - INTERVAL '23 hours');

-- OAUTH USERS
INSERT INTO oauth_user (user_id, provider_name, provider_user_id, access_token, refresh_token) VALUES
                                                                                                   (1, 'google',   'google_uid_budi', 'access_budi_google', 'refresh_budi_google'),
                                                                                                   (2, 'facebook', 'fb_uid_siti',     'access_siti_fb',     'refresh_siti_fb'),
                                                                                                   (3, 'google',   'google_uid_andi', 'access_andi_google', 'refresh_andi_google');

-- RESET SEQUENCES
SELECT setval('users_id_seq',                   (SELECT MAX(id) FROM users));
SELECT setval('category_payment_method_id_seq', (SELECT MAX(id) FROM category_payment_method));
SELECT setval('payment_method_id_seq',          (SELECT MAX(id) FROM payment_method));
SELECT setval('transactions_id_seq',            (SELECT MAX(id) FROM transactions));