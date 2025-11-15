import requests
import json
import random
import string
import time

#__Codigo generado con Inteligencia Artificial Generativa de OpenAI__

# Deshabilitar los warnings de certificados autofirmados
requests.packages.urllib3.disable_warnings()

# MODIFICADO: Apuntamos a localhost
BASE_URL = "https://myserver.local:5000"

def random_string(length=10):
    """Genera una cadena aleatoria."""
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(length))

def test_version():
    print("--- Test GET /version ---")
    response = requests.get(f"{BASE_URL}/version", verify=False)
    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.text}")
    assert response.status_code == 200
    assert "Database as a Service" in response.text
    print("✅ OK\n")

# ... (el resto del fichero de tests no necesita cambios)
def test_auth_flow():
    print("--- Test Flujo de Autenticación ---")
    username = random_string()
    password = "testpassword123"

    # 1. Signup
    print("Probando POST /signup...")
    signup_data = {"username": username, "password": password}
    response = requests.post(f"{BASE_URL}/signup", json=signup_data, verify=False)
    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.json()}")
    assert response.status_code == 200
    assert "access_token" in response.json()
    token = response.json()["access_token"]
    print("✅ Signup OK\n")

    # 2. Signup con usuario existente
    print("Probando POST /signup con usuario existente...")
    response = requests.post(f"{BASE_URL}/signup", json=signup_data, verify=False)
    print(f"Status Code: {response.status_code}")
    assert response.status_code == 409
    print("✅ Conflicto de usuario OK\n")

    # 3. Login
    print("Probando POST /login...")
    login_data = {"username": username, "password": password}
    response = requests.post(f"{BASE_URL}/login", json=login_data, verify=False)
    print(f"Status Code: {response.status_code}")
    assert response.status_code == 200
    assert "access_token" in response.json()
    print("✅ Login OK\n")

    return username, token

def test_documents_flow(username, token):
    print("--- Test Flujo de Documentos ---")
    headers = {"Authorization": f"token {token}"}
    doc_id = "mi-primer-documento"
    doc_content = {"titulo": "Hola Mundo", "contenido": "Este es mi primer doc."}

    # 1. Crear documento
    print(f"Probando POST /{username}/{doc_id}...")
    create_data = {"doc_content": doc_content}
    response = requests.post(f"{BASE_URL}/{username}/{doc_id}", json=create_data, headers=headers, verify=False)
    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.json()}")
    assert response.status_code == 200
    assert "size" in response.json()
    print("✅ Creación de documento OK\n")

    # 2. Obtener documento
    print(f"Probando GET /{username}/{doc_id}...")
    response = requests.get(f"{BASE_URL}/{username}/{doc_id}", headers=headers, verify=False)
    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.json()}")
    assert response.status_code == 200
    assert response.json() == doc_content
    print("✅ Obtención de documento OK\n")

    # 3. Actualizar documento
    print(f"Probando PUT /{username}/{doc_id}...")
    updated_content = {"titulo": "Hola Mundo v2", "contenido": "Contenido actualizado."}
    update_data = {"doc_content": updated_content}
    response = requests.put(f"{BASE_URL}/{username}/{doc_id}", json=update_data, headers=headers, verify=False)
    print(f"Status Code: {response.status_code}")
    assert response.status_code == 200
    print("✅ Actualización de documento OK\n")

    # Verificar actualización
    response = requests.get(f"{BASE_URL}/{username}/{doc_id}", headers=headers, verify=False)
    assert response.json() == updated_content

    # 4. Obtener todos los documentos
    print(f"Probando GET /{username}/_all_docs...")
    response = requests.get(f"{BASE_URL}/{username}/_all_docs", headers=headers, verify=False)
    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.json()}")
    assert response.status_code == 200
    assert doc_id in response.json()
    assert response.json()[doc_id] == updated_content
    print("✅ Obtener todos los documentos OK\n")

    # 5. Borrar documento
    print(f"Probando DELETE /{username}/{doc_id}...")
    response = requests.delete(f"{BASE_URL}/{username}/{doc_id}", headers=headers, verify=False)
    print(f"Status Code: {response.status_code}")
    assert response.status_code == 200
    print("✅ Borrado de documento OK\n")

    # Verificar borrado
    response = requests.get(f"{BASE_URL}/{username}/{doc_id}", headers=headers, verify=False)
    assert response.status_code == 404

def test_security_checks():
    print("--- Test de Seguridad ---")
    # Intentar acceder a un recurso sin token
    print("Probando acceso sin token...")
    response = requests.get(f"{BASE_URL}/testuser/some_doc", verify=False)
    print(f"Status Code: {response.status_code}")
    assert response.status_code == 401
    print("✅ Acceso denegado sin token OK\n")

    # Intentar acceder al recurso de otro usuario
    print("Probando acceso a recurso de otro usuario...")
    user1, token1 = "user" + random_string(5), ""
    user2, token2 = "user" + random_string(5), ""

    # Crear usuarios y obtener tokens
    requests.post(f"{BASE_URL}/signup", json={"username": user1, "password": "pw1"}, verify=False)
    resp = requests.post(f"{BASE_URL}/login", json={"username": user1, "password": "pw1"}, verify=False)
    token1 = resp.json()["access_token"]

    requests.post(f"{BASE_URL}/signup", json={"username": user2, "password": "pw2"}, verify=False)
    resp = requests.post(f"{BASE_URL}/login", json={"username": user2, "password": "pw2"}, verify=False)
    token2 = resp.json()["access_token"]
    
    # User1 crea un documento
    doc_id = "doc_de_user1"
    doc_content = {"owner": user1}
    requests.post(f"{BASE_URL}/{user1}/{doc_id}", json={"doc_content": doc_content}, headers={"Authorization": f"token {token1}"}, verify=False)

    # User2 intenta leer el documento de User1
    response = requests.get(f"{BASE_URL}/{user1}/{doc_id}", headers={"Authorization": f"token {token2}"}, verify=False)
    print(f"Status Code: {response.status_code}")
    assert response.status_code == 403
    print("✅ Acceso denegado a recurso ajeno OK\n")


if __name__ == "__main__":
    test_version()
    username, token = test_auth_flow()
    test_documents_flow(username, token)
    test_security_checks()
    print("🎉 ¡Todos los tests han pasado! 🎉")