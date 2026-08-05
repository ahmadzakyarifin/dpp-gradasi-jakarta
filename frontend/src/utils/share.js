/**
 * Utility helper for modern web sharing.
 * Uses Web Share API (native OS share) when available,
 * and provides fallbacks for specific platforms (WhatsApp, Facebook, Twitter, LinkedIn)
 * or copying to clipboard.
 */

export async function shareContent({ title, text, url }) {
  const shareUrl = url || window.location.href;
  const shareData = {
    title: title,
    text: text || title,
    url: shareUrl,
  };

  if (navigator.share && navigator.canShare && navigator.canShare(shareData)) {
    try {
      await navigator.share(shareData);
      return { success: true, method: 'native' };
    } catch (err) {
      if (err.name !== 'AbortError') {
        console.error('Error with navigator.share:', err);
      } else {
        return { success: false, reason: 'cancelled' };
      }
    }
  }
  return { success: false, method: 'fallback' };
}

export function getShareUrl(platform, { title, text, url }) {
  const shareUrl = url || window.location.href;
  const cleanText = text ? text.replace(/<[^>]*>/g, '').trim() : ''; // strip HTML if any
  const formattedText = `*${title}*\n\n${cleanText ? cleanText + '\n\n' : ''}Baca selengkapnya di: ${shareUrl}`;
  
  const encodedUrl = encodeURIComponent(shareUrl);
  const encodedText = encodeURIComponent(formattedText);
  const encodedTitle = encodeURIComponent(title);

  switch (platform) {
    case 'whatsapp':
      return `https://api.whatsapp.com/send?text=${encodedText}`;
    case 'facebook':
      // Facebook sharer only accepts url. Description is scraped from Open Graph meta tags.
      return `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}`;
    case 'twitter':
    case 'x':
      return `https://twitter.com/intent/tweet?text=${encodedTitle}&url=${encodedUrl}`;
    case 'linkedin':
      return `https://www.linkedin.com/sharing/share-offsite/?url=${encodedUrl}`;
    default:
      return '';
  }
}

export async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch (err) {
    console.error('Failed to copy to clipboard', err);
    return false;
  }
}
