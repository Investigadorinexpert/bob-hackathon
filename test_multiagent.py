#!/usr/bin/env python3
"""
Script para probar el sistema multiagente de BOB Chatbot
"""
import requests
import json
import time

BASE_URL = "http://localhost:3000/api"
session_id = None

def send_message(message, channel="whatsapp"):
    global session_id

    payload = {
        "message": message,
        "channel": channel
    }

    if session_id:
        payload["sessionId"] = session_id

    print(f"\n{'='*80}")
    print(f"👤 USER: {message}")
    print(f"{'='*80}")

    try:
        response = requests.post(f"{BASE_URL}/chat/message", json=payload, timeout=30)
        response.raise_for_status()

        data = response.json()
        session_id = data.get("sessionId")

        print(f"🤖 BOT: {data.get('reply', 'No response')}")
        print(f"\n📊 Score: {data.get('leadScore', 0)}/100 | Category: {data.get('category', 'N/A')}")
        print(f"🆔 Session: {session_id}")

        return data

    except requests.exceptions.RequestException as e:
        print(f"❌ Error: {e}")
        return None

def test_faq_routing():
    """Test 1: Verificar routing a FAQ Agent"""
    print("\n" + "="*80)
    print("TEST 1: FAQ ROUTING - Preguntas sobre el proceso")
    print("="*80)

    send_message("Hola, ¿cómo funciona el proceso de subasta?")
    time.sleep(2)

def test_auction_routing():
    """Test 2: Verificar routing a Auction Agent"""
    print("\n" + "="*80)
    print("TEST 2: AUCTION ROUTING - Búsqueda de vehículos")
    print("="*80)

    send_message("Busco una camioneta Toyota para mi negocio")
    time.sleep(2)

def test_scoring_hot_lead():
    """Test 3: Crear un HOT lead con conversación completa"""
    print("\n" + "="*80)
    print("TEST 3: SCORING - HOT LEAD (85-100 puntos)")
    print("="*80)

    # Reset session
    global session_id
    session_id = None

    # Conversación de un lead caliente
    messages = [
        "Hola, soy empresario de transporte en Lima",
        "Necesito 2 camionetas 4x4 urgente para esta semana",
        "Tengo presupuesto de S/100,000 por unidad",
        "Ya he comprado en subastas antes, conozco el proceso",
        "Quiero ver Toyota Hilux o similar, año 2020 en adelante",
        "¿Pueden hacer una inspección técnica? Necesito garantías",
        "Perfecto, estoy disponible mañana para coordinar"
    ]

    for msg in messages:
        send_message(msg)
        time.sleep(3)  # Simular conversación natural

def test_spam_detection():
    """Test 4: Verificar detección de spam"""
    print("\n" + "="*80)
    print("TEST 4: SPAM DETECTION")
    print("="*80)

    global session_id
    session_id = None

    send_message("COMPRA AHORA!!! OFERTA INCREIBLE!!! CLICK AQUI!!!")
    time.sleep(2)

def test_ambiguous_message():
    """Test 5: Verificar manejo de mensajes ambiguos"""
    print("\n" + "="*80)
    print("TEST 5: AMBIGUOUS MESSAGE HANDLING")
    print("="*80)

    global session_id
    session_id = None

    send_message("hola")
    time.sleep(2)

def main():
    print("\n" + "🚀"*40)
    print("PRUEBA DEL SISTEMA MULTIAGENTE - BOB CHATBOT")
    print("🚀"*40)

    # Verificar que el servidor esté corriendo
    try:
        response = requests.get("http://localhost:3000/health", timeout=5)
        print(f"\n✅ Backend online: {response.json()}")
    except:
        print("\n❌ Backend no está corriendo. Ejecuta: go run cmd/server/main.go")
        return

    # Ejecutar tests
    test_faq_routing()
    test_auction_routing()
    test_spam_detection()
    test_ambiguous_message()
    test_scoring_hot_lead()

    print("\n" + "="*80)
    print("✅ PRUEBAS COMPLETADAS")
    print("="*80)

if __name__ == "__main__":
    main()
