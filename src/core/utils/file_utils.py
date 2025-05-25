"""
File utilities for the Discord music bot
Helper functions for file operations
"""

import os
import json
from pathlib import Path

def get_project_root():
    """Get the absolute path to the project root directory."""
    # Go up from src/core/utils to project root
    return Path(__file__).parents[3]

def read_env_file(env_path=None):
    """
    Read environment variables from a .env file.
    
    Args:
        env_path: Path to the .env file. If None, use the one in project root.
        
    Returns:
        dict: Environment variables
    """
    if env_path is None:
        env_path = os.path.join(get_project_root(), '.env')
        
    if not os.path.exists(env_path):
        return {}
        
    env_vars = {}
    with open(env_path, 'r') as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith('#') or '//' in line:
                continue
                
            if '=' in line:
                key, value = line.split('=', 1)
                env_vars[key.strip()] = value.strip()
                
    return env_vars

def write_env_file(env_vars, env_path=None):
    """
    Write environment variables to a .env file.
    
    Args:
        env_vars: Dictionary of environment variables
        env_path: Path to the .env file. If None, use the one in project root.
    """
    if env_path is None:
        env_path = os.path.join(get_project_root(), '.env')
        
    with open(env_path, 'w') as f:
        f.write("# Environment variables for the Discord music bot\n")
        for key, value in env_vars.items():
            f.write(f"{key}={value}\n")

def write_config(config, filename='config.json'):
    """
    Write configuration to a JSON file.
    
    Args:
        config: Configuration dictionary
        filename: Name of the config file
    """
    config_path = os.path.join(get_project_root(), filename)
    with open(config_path, 'w') as f:
        json.dump(config, f, indent=2)
        
def read_config(filename='config.json'):
    """
    Read configuration from a JSON file.
    
    Args:
        filename: Name of the config file
        
    Returns:
        dict: Configuration
    """
    config_path = os.path.join(get_project_root(), filename)
    if not os.path.exists(config_path):
        return {}
        
    with open(config_path, 'r') as f:
        return json.load(f)

def ensure_directory_exists(directory):
    """
    Ensure that a directory exists, creating it if it doesn't.
    
    Args:
        directory: Directory path
    """
    if not os.path.exists(directory):
        os.makedirs(directory)
        
def find_ffmpeg():
    """
    Find the FFmpeg executable on the system.
    
    Returns:
        str: Path to FFmpeg executable, or None if not found
    """
    # Common locations to check
    common_paths = [
        '/usr/bin/ffmpeg',
        '/usr/local/bin/ffmpeg',
        '/opt/homebrew/bin/ffmpeg',  # Common on macOS with Homebrew
        'ffmpeg'  # Just the command name, hoping it's in PATH
    ]
    
    # Check each path
    for path in common_paths:
        try:
            if path != 'ffmpeg':  # Skip checking file existence for just the command name
                if os.path.exists(path):
                    return path
            else:
                import subprocess
                subprocess.run([path, '-version'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                return path
        except:
            continue
    
    return None