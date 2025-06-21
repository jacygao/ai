import time
from fastapi import FastAPI, Request, Form
from fastapi.templating import Jinja2Templates
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse
import os

from app.models import SearchRequest, SearchResponse, SearchMode
from app.search.keyword_search import KeywordSearch
from app.search.vector_search import VectorSearch
from app.utils.data_loader import data_loader

# Initialize FastAPI app
app = FastAPI(
    title="AI Search vs Keyword Search Demo",
    description="A demonstration of AI-powered vector search vs traditional keyword search",
    version="1.0.0"
)

# Mount static files
app.mount("/static", StaticFiles(directory="static"), name="static")

# Setup templates
templates = Jinja2Templates(directory="templates")

# Initialize search engines
keyword_search = KeywordSearch()
vector_search = VectorSearch()

@app.get("/", response_class=HTMLResponse)
async def home(request: Request):
    """Main application page"""
    return templates.TemplateResponse(
        "index.html",
        {
            "request": request,
            "document_count": data_loader.get_document_count()
        }
    )

@app.post("/api/search")
async def search(request: SearchRequest) -> SearchResponse:
    """Search endpoint for both keyword and vector search"""
    start_time = time.time()
    
    try:
        if request.mode == SearchMode.KEYWORD:
            results = keyword_search.search(request.query, request.limit)
        else:  # vector search
            results = vector_search.search(request.query, request.limit)
        
        search_time = (time.time() - start_time) * 1000  # Convert to milliseconds
        
        return SearchResponse(
            results=results,
            total_results=len(results),
            query=request.query,
            mode=request.mode,
            search_time_ms=search_time
        )
    
    except Exception as e:
        print(f"Search error: {e}")
        return SearchResponse(
            results=[],
            total_results=0,
            query=request.query,
            mode=request.mode,
            search_time_ms=0.0
        )

@app.get("/api/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "document_count": data_loader.get_document_count(),
        "keyword_search_ready": keyword_search.bm25 is not None,
        "vector_search_ready": vector_search.model is not None
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000) 