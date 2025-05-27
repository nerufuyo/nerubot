"""
Logging utilities for the Discord music bot
This file contains functions to set up and use logging consistently across the app
"""

import logging
import sys
import os
from pathlib import Path

def setup_logger(name=None, level=logging.INFO, log_to_file=True, log_file='bot.log'):
    """
    Set up a logger with consistent formatting.
    
    Args:
        name: Logger name (usually __name__ from the calling module)
        level: Logging level
        log_to_file: Whether to log to file
        log_file: Path to log file
        
    Returns:
        Logger instance
    """
    logger = logging.getLogger(name)
    logger.setLevel(level)
    
    # Clear existing handlers if any
    logger.handlers = []
    
    # Create formatter
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    
    # Add console handler
    console_handler = logging.StreamHandler(stream=sys.stdout)
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)
    
    # Add file handler if requested
    if log_to_file:
        # Determine log file path - use absolute path from project root
        if not os.path.isabs(log_file):
            # Go up from src/core/utils to project root
            project_root = Path(__file__).parents[3]
            log_file = os.path.join(project_root, log_file)
            
        file_handler = logging.FileHandler(filename=log_file, mode='a')
        file_handler.setFormatter(formatter)
        logger.addHandler(file_handler)
    
    return logger

def get_logger(name=None):
    """
    Get an existing logger or create a new one if it doesn't exist.
    
    Args:
        name: Logger name
        
    Returns:
        Logger instance
    """
    # Just return the named logger without reconfiguring
    # This avoids duplicate configuration
    return logging.getLogger(name)

def configure_root_logger(level=logging.DEBUG, log_file='bot.log'):
    """
    Configure the root logger for the application.
    This should be called once at application startup.
    
    Args:
        level: Logging level
        log_file: Path to log file
    """
    # Get the root logger
    root_logger = logging.getLogger()
    
    # If the root logger already has handlers, don't add more
    if root_logger.handlers:
        return root_logger
    
    # Set the root logger's level
    root_logger.setLevel(level)
    
    # Create formatter
    formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    
    # Add console handler
    console_handler = logging.StreamHandler(stream=sys.stdout)
    console_handler.setFormatter(formatter)
    root_logger.addHandler(console_handler)
    
    # Add file handler
    file_handler = logging.FileHandler(filename=log_file, mode='a')
    file_handler.setFormatter(formatter)
    root_logger.addHandler(file_handler)
    
    # Set levels for some noisy libraries
    logging.getLogger('discord').setLevel(logging.WARNING)
    logging.getLogger('discord.http').setLevel(logging.WARNING)
    logging.getLogger('websockets').setLevel(logging.WARNING)
    logging.getLogger('asyncio').setLevel(logging.WARNING)
    
    # Return the root logger
    return root_logger