from diagrams import Diagram
from diagrams.onprem.client import Client
from diagrams.onprem.compute import Server
from diagrams.onprem.queue import Kafka
from diagrams.onprem.database import PostgreSQL

with Diagram("Chat Application Architecture (On-Prem)", show=True, direction="LR") as diag:

    # Producer service with custom color and style
    producer = Server("Producer Service")
    producer.style = "filled,rounded,fillcolor=skyblue, fontcolor=black"

    # Kafka Broker with custom color and style
    kafka = Kafka("Kafka Broker")
    kafka.style = "filled,rounded,fillcolor=orange, fontcolor=black"

    # Consumer Service with custom color and style
    consumer = Server("Consumer Service")
    consumer.style = "filled,rounded,fillcolor=lightgreen, fontcolor=black"

    # Database with custom color and style
    db = PostgreSQL("Database")
    db.style = "filled,rounded,fillcolor=lightyellow, fontcolor=black"

    # WebSocket clients with custom color and style
    client = Client("WebSocket Clients")
    client.style = "filled,rounded,fillcolor=lightcoral, fontcolor=black"

    # Additional services with colors
    # Monitoring service
    monitoring = Server("Monitoring Service")
    monitoring.style = "filled,rounded,fillcolor=lightgray, fontcolor=black"

  

    # Connect all components with arrows
    producer >> kafka >> consumer
    consumer >> db
    consumer >> client
    consumer >> monitoring  # Connect consumer to monitoring
    

diag