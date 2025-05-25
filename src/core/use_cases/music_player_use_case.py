"""
Music player use cases for the Discord music bot
"""
from collections import deque
from typing import Optional, List, Tuple
from src.core.entities.song import Song
import random


class MusicPlayerUseCase:
    """Use cases for the music player."""
    
    def __init__(self):
        self.queues = {}
        self.loop_modes = {}  # guild_id -> 'off', 'song', or 'queue'
    
    def add_song(self, guild_id: int, song: Song) -> None:
        """
        Add a song to the queue.
        
        Args:
            guild_id: The Discord guild ID
            song: The song to add
        """
        if guild_id not in self.queues:
            self.queues[guild_id] = deque()
        self.queues[guild_id].append(song)
    
    def get_current_song(self, guild_id: int) -> Optional[Song]:
        """
        Get the current song for a guild.
        
        Args:
            guild_id: The Discord guild ID
            
        Returns:
            Optional[Song]: The current song or None
        """
        queue = self.queues.get(guild_id)
        if not queue:
            return None
        return queue[0] if queue else None
    
    def skip(self, guild_id: int) -> Optional[Song]:
        """
        Skip to the next song.
        
        Args:
            guild_id: The Discord guild ID
            
        Returns:
            Optional[Song]: The next song or None
        """
        queue = self.queues.get(guild_id)
        if not queue:
            return None
            
        loop_mode = self.loop_modes.get(guild_id, 'off')
        
        if loop_mode == 'song':
            return queue[0]
        elif loop_mode == 'queue':
            # Move current song to the end and return the next one
            current = queue.popleft()
            queue.append(current)
        else:
            # Remove current song
            queue.popleft()
            
        return queue[0] if queue else None
    
    def clear_queue(self, guild_id: int) -> None:
        """
        Clear the queue for a guild.
        
        Args:
            guild_id: The Discord guild ID
        """
        if guild_id in self.queues:
            self.queues[guild_id].clear()
    
    def get_queue(self, guild_id: int, start: int = 0, count: int = 10) -> Tuple[List[Song], int]:
        """
        Get songs from the queue.
        
        Args:
            guild_id: The Discord guild ID
            start: The starting index
            count: The number of songs to get
            
        Returns:
            Tuple[List[Song], int]: The songs and the total number of songs
        """
        queue = self.queues.get(guild_id, deque())
        songs = list(queue)[start:start + count]
        return songs, len(queue)
    
    def shuffle(self, guild_id: int) -> None:
        """
        Shuffle the queue for a guild.
        
        Args:
            guild_id: The Discord guild ID
        """
        queue = self.queues.get(guild_id)
        if queue and len(queue) > 1:
            current = queue.popleft()  # Keep current song
            remaining = list(queue)
            random.shuffle(remaining)
            queue.clear()
            queue.append(current)  # Restore current song
            queue.extend(remaining)
    
    def set_loop_mode(self, guild_id: int, mode: str) -> None:
        """
        Set the loop mode for a guild.
        
        Args:
            guild_id: The Discord guild ID
            mode: The loop mode ('off', 'song', or 'queue')
        """
        if mode in ['off', 'song', 'queue']:
            self.loop_modes[guild_id] = mode
