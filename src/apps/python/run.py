#!/usr/bin/env python3
"""
AI Search vs Keyword Search Demo
Application runner script
"""

import uvicorn
import os
import sys

def main():
    """Main application entry point"""
    print("🚀 Starting AI Search vs Keyword Search Demo...")
    print("📚 Loading documents and initializing search engines...")
    
    # Add the current directory to Python path
    sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
    
    # Import the app after setting up the path
    from app.main import app
    
    print("✅ Application initialized successfully!")
    print("🌐 Starting server at http://localhost:8000")
    print("📖 Open your browser and navigate to the URL above")
    print("🔄 Press Ctrl+C to stop the server")
    
    # Run the application
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=8000,
        reload=False,  # Set to True for development
        log_level="info"
    )

if __name__ == "__main__":
    main() 