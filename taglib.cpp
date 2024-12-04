//go:build ignore
#include <cstring>
#include <iostream>

#include "fileref.h"
#include "tpropertymap.h"

__attribute__((export_name("malloc"))) void *exported_malloc(size_t size) {
  return malloc(size);
}

__attribute__((export_name("taglib_file_tags"))) char **
taglib_file_tags(const char *filename) {
  TagLib::FileRef file(filename);
  if (file.isNull())
    return nullptr;

  auto properties = file.properties();

  size_t len = 0;
  for (const auto &kvs : properties)
    len += kvs.second.size();

  char **tags = static_cast<char **>(malloc(sizeof(char *) * (len + 1)));
  if (!tags)
    return nullptr;

  size_t i = 0;
  for (const auto &kvs : properties)
    for (const auto &v : kvs.second) {
      TagLib::String row;
      row.append(kvs.first);
      row.append("\t");
      row.append(v);
      tags[i] = static_cast<char *>(malloc(row.size() + 1));
      strncpy(tags[i], row.toCString(), row.size());
      tags[i][row.size()] = '\0';
      i++;
    }
  tags[len] = nullptr;

  return tags;
}

__attribute__((export_name("taglib_file_write_tags"))) bool
taglib_file_write_tags(const char *filename, const char **tags) {
  if (!filename || !tags)
    return false;

  TagLib::FileRef file(filename);
  if (file.isNull())
    return false;

  TagLib::PropertyMap properties;
  for (size_t i = 0; tags[i] != NULL; i++) {
    TagLib::String row = tags[i];
    if (auto ti = row.find("\t"); ti >= 0) {
      TagLib::String key(row.substr(0, ti));
      TagLib::StringList value(row.substr(ti + 1));
      properties.insert(key, value);
    }
  }

  if (auto rejected = file.setProperties(properties); rejected.size() > 0)
    return 0;

  return file.save();
}

__attribute__((export_name("taglib_file_audioproperties"))) int *
taglib_file_audioproperties(const char *filename) {
  TagLib::FileRef file(filename);

  int *arr = static_cast<int *>(malloc(4 * sizeof(int)));

  auto audioProperties = file.audioProperties();
  arr[0] = audioProperties->lengthInMilliseconds();
  arr[1] = audioProperties->channels();
  arr[2] = audioProperties->sampleRate();
  arr[3] = audioProperties->bitrate();

  return arr;
}